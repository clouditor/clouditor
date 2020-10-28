package io.clouditor.data_access_layer;

import io.clouditor.assurance.*;
import io.clouditor.assurance.ccl.*;
import io.clouditor.auth.User;
import io.clouditor.discovery.Asset;
import io.clouditor.discovery.AssetProperties;
import io.clouditor.discovery.DiscoveryResult;
import io.clouditor.discovery.Scan;
import java.io.Serializable;
import java.util.*;
import java.util.function.Consumer;
import java.util.function.Function;
import javax.persistence.*;
import javax.persistence.criteria.CriteriaQuery;
import org.hibernate.HibernateException;
import org.hibernate.Session;
import org.hibernate.SessionFactory;
import org.hibernate.Transaction;
import org.hibernate.cfg.Configuration;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class HibernatePersistence {

  private static final String USER_NAME = "postgres";
  private static final String PASSWORD = "postgres";
  private static final String DB_NAME = "postgres";

  private static final Configuration CONFIGURATION =
      new Configuration()
          .configure()
          .addAnnotatedClass(Domain.class)
          .addAnnotatedClass(Certification.class)
          .addAnnotatedClass(Control.class)
          .addAnnotatedClass(Rule.class)
          .addAnnotatedClass(EvaluationResult.class)
          .addAnnotatedClass(AssetType.class)
          .addAnnotatedClass(FilteredAssetType.class)
          .addAnnotatedClass(Asset.class)
          .addAnnotatedClass(Condition.class)
          .addAnnotatedClass(Scan.class)
          .addAnnotatedClass(DiscoveryResult.class)
          .addAnnotatedClass(User.class)
          .setProperty("hibernate.connection.username", USER_NAME)
          .setProperty("hibernate.connection.password", PASSWORD)
          .setProperty("hibernate.enable_lazy_load_no_trans", "true")
          .setProperty("hibernate.dialect", "org.hibernate.dialect.PostgreSQL94Dialect")
          .setProperty("hibernate.hbm2ddl.auto", "create-drop");

  private static SessionFactory sessionFactory;

  protected static final Logger LOGGER = LoggerFactory.getLogger(HibernatePersistence.class);

  private static SessionFactory getSessionFactory() {
    return sessionFactory;
  }

  public static void init() {
    init("localhost", 5432, DB_NAME);
  }

  public static void init(final String host, final int port, final String dbName) {
    Objects.requireNonNull(host);
    Objects.requireNonNull(dbName);
    if (port < 1024 || port > 60000) throw new IllegalArgumentException();
    sessionFactory =
        CONFIGURATION
            .setProperty("hibernate.connection.driver_class", "org.postgresql.Driver")
            .setProperty(
                "hibernate.connection.url", "jdbc:postgresql://" + host + ":" + port + "/" + dbName)
            .buildSessionFactory();
  }

  public static void init(final String dbName) {
    Objects.requireNonNull(dbName);
    sessionFactory =
        CONFIGURATION
            .setProperty("hibernate.connection.driver_class", "org.h2.Driver")
            .setProperty("hibernate.connection.url", "jdbc:h2:~/" + dbName)
            .buildSessionFactory();
  }

  private <T> Optional<T> exec(final Function<Session, T> function) {
    Objects.requireNonNull(function);
    if (getSessionFactory() == null)
      throw new IllegalStateException("The Database Connection is not initialized.");
    Optional<T> result = Optional.empty();
    Optional<Transaction> transaction = Optional.empty();
    try (final Session session = sessionFactory.openSession()) {
      transaction = Optional.of(session.beginTransaction());
      result = Optional.ofNullable(function.apply(session));
      transaction.ifPresent(EntityTransaction::commit);
    } catch (final HibernateException exception) {
      exception.printStackTrace();
      transaction.ifPresent(
          t -> {
            if (t.getStatus().canRollback()) t.rollback();
          });
    }
    return result;
  }

  private void execConsumer(final Consumer<Session> consumer) {
    Objects.requireNonNull(consumer);
    exec(
        session -> {
          consumer.accept(session);
          return new Object();
        });
  }

  public <T> void saveOrUpdate(final T toSave) {
    execConsumer(session -> session.saveOrUpdate(toSave));
  }

  public <T> Optional<T> get(final Class<T> resultType, final Serializable primaryKey) {
    return exec(session -> session.get(resultType, primaryKey));
  }

  public <T> List<T> listAll(final Class<T> resultType) {
    return exec(session -> {
          final CriteriaQuery<T> criteriaQuery =
              session.getCriteriaBuilder().createQuery(resultType);
          criteriaQuery.from(resultType);
          return session.createQuery(criteriaQuery).getResultList();
        })
        .orElseThrow();
  }

  public <T> void delete(final T toDelete) {
    execConsumer(session -> session.delete(toDelete));
  }

  public <T> void delete(final Class<T> deleteType, final Serializable id) {
    execConsumer(session -> session.delete(session.load(deleteType, id)));
  }

  public <T> int count(final Class<T> countType) {
    return listAll(countType).size();
  }

  public static void close() {
    if (sessionFactory != null) sessionFactory.close();
  }

  public static void main(final String... args) {
    init(DB_NAME);

    final HibernatePersistence persistence = new HibernatePersistence();

    final String sutPK = "username";
    final User sut = new User(sutPK, "password");
    sut.setFullName("fullName");
    sut.setEmail("username@test.edu");
    sut.setShadow(true);
    final String tableName = "c_user";

    final String assetTypeID = "asset_type_id";
    final AssetType assetType = new AssetType();
    assetType.setValue(assetTypeID);
    final String assetTypeTableName = "asset_type";

    final String filteredAssetTypeID = "filtered_asset_type_ID";
    final FilteredAssetType filteredAssetType = new FilteredAssetType();
    filteredAssetType.setValue(filteredAssetTypeID);
    final String filteredAssetTypTableName = "filtered_asset_type";

    final Condition condition = new Condition();
    condition.setAssetType(assetType);
    condition.setSource("source");
    final Condition.ConditionPK conditionID = condition.getConditionPK();
    final String conditionTableName = "condition";

    final String domainID = "domain_name";
    final Domain domain = new Domain(domainID);
    domain.setDescription("domain description");
    final String domainTableName = "cloud_domain";

    final String controlID = "control_id";
    final Control control = new Control();
    control.setControlId(controlID);
    control.setDescription("Control Description");
    control.setFulfilled(Control.Fulfillment.GOOD);
    control.setDomain(domain);
    control.setName("control name");
    control.setAutomated(true);
    control.setActive(true);
    final String controlTableName = "control";

    final String ruleID = "rule_id";
    final Rule rule = new Rule();
    rule.setId(ruleID);
    rule.setActive(true);
    rule.setName("rule name");
    rule.setDescription("rule description");
    rule.getControls().add(control);
    rule.setCondition(condition);
    final String ruleTableName = "rule";

    final String certificationID = "certification_id";
    final Certification certification = new Certification();
    certification.setId(certificationID);
    certification.setDescription("Certification Description");
    certification.setPublisher("Certification Publisher");
    certification.setWebsite("website");
    certification.setControls(List.of(control));
    final String certificationTableName = "certification";

    final AssetProperties assetProperties = new AssetProperties("TEST_KEY", "TEST_VALUE");
    final EvaluationResult evaluationResult = new EvaluationResult(rule, assetProperties);
    final String evaluationResultTableName = "evaluation_result";
    final String evaluationResultID = evaluationResult.getTimeStamp();

    /*
    persistence.saveOrUpdate(
            sut,
            domain,
            control,
            certification,
            assetType,
            filteredAssetType,
            condition,
            rule,
            evaluationResult
        );
    */

    test(
        persistence,
        sut,
        User.class,
        sutPK,
        tableName,
        () ->
            test(
                persistence,
                assetType,
                AssetType.class,
                assetTypeID,
                assetTypeTableName,
                () ->
                    test(
                        persistence,
                        filteredAssetType,
                        FilteredAssetType.class,
                        filteredAssetTypeID,
                        filteredAssetTypTableName,
                        () ->
                            test(
                                persistence,
                                condition,
                                Condition.class,
                                conditionID,
                                conditionTableName,
                                () ->
                                    test(
                                        persistence,
                                        domain,
                                        Domain.class,
                                        domainID,
                                        domainTableName,
                                        () ->
                                            test(
                                                persistence,
                                                control,
                                                Control.class,
                                                controlID,
                                                controlTableName,
                                                () ->
                                                    test(
                                                        persistence,
                                                        rule,
                                                        Rule.class,
                                                        ruleID,
                                                        ruleTableName,
                                                        () ->
                                                            test(
                                                                persistence,
                                                                certification,
                                                                Certification.class,
                                                                certificationID,
                                                                certificationTableName,
                                                                () -> {
                                                                  control
                                                                      .getResults()
                                                                      .add(evaluationResult);
                                                                  persistence.saveOrUpdate(control);
                                                                  test(
                                                                      persistence,
                                                                      evaluationResult,
                                                                      EvaluationResult.class,
                                                                      evaluationResultID,
                                                                      evaluationResultTableName,
                                                                      () -> {});
                                                                }))))))));
    sessionFactory.close();
  }

  private static <T> void test(
      final HibernatePersistence persistence,
      final T sut,
      final Class<T> type,
      final Serializable sutPK,
      final String tableName,
      final Runnable runnable) {
    LOGGER.info(tableName + " " + persistence.listAll(type));

    boolean isPresent = persistence.get(type, sutPK).isPresent();
    LOGGER.info(tableName + " IS_PRESENT: " + isPresent);

    if (!isPresent) persistence.saveOrUpdate(sut);

    LOGGER.info(tableName + " " + persistence.listAll(type));

    var storedValue = persistence.get(type, sutPK);
    LOGGER.info(tableName + " STORED_VALUE: " + storedValue);

    runnable.run();

    storedValue.ifPresent(persistence::delete);

    isPresent = persistence.get(type, sutPK).isPresent();
    LOGGER.info(tableName + " IS_PRESENT: " + isPresent);

    LOGGER.info(tableName + " " + persistence.listAll(type));
  }
}
