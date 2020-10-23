package io.clouditor.data_access_layer;

import io.clouditor.assurance.*;
import io.clouditor.assurance.ccl.AssetType;
import io.clouditor.assurance.ccl.Condition;
import io.clouditor.assurance.ccl.FilteredAssetType;
import io.clouditor.auth.User;
import io.clouditor.data_access_layer.utils.BoomingFunction;
import io.clouditor.data_access_layer.utils.DBUtils;
import io.clouditor.discovery.Asset;
import io.clouditor.discovery.DiscoveryResult;
import io.clouditor.discovery.Scan;
import org.hibernate.HibernateException;
import org.hibernate.Session;
import org.hibernate.SessionFactory;
import org.hibernate.Transaction;
import org.hibernate.cfg.Configuration;

import javax.persistence.*;
import java.io.Closeable;
import java.io.Serializable;
import java.sql.SQLException;
import java.util.*;
import java.util.function.Consumer;

public class Persistence implements Closeable {

    private final SessionFactory factory = new Configuration()
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
            .buildSessionFactory();

    private SessionFactory getFactory() {
        return factory;
    }

    private <T> Optional<T> exec(final BoomingFunction<Session, T> function) {
        Objects.requireNonNull(function);
        Optional<T> result = Optional.empty();
        Optional<Transaction> transaction = Optional.empty();
        try (final Session session = getFactory().openSession()) {
            transaction = Optional.of(session.beginTransaction());
            result = Optional.ofNullable(
                    function.handleSQLException()
                            .apply(session)
            );
            transaction.ifPresent(EntityTransaction::commit);
        } catch (final HibernateException exception) {
            exception.printStackTrace();
            transaction.ifPresent(t -> {
                if (t.getStatus().canRollback())
                    t.rollback();
            });
        }
        return result;
    }

    private void execConsumer(final Consumer<Session> consumer) {
        Objects.requireNonNull(consumer);
        exec(session -> {
           consumer.accept(session);
           return new Object();
        });
    }

    public <T> void save(final T toSave) {
        exec(session -> session.save(toSave));
    }

    public <T> Optional<T> get(final Class<T> resultType, final Serializable primaryKey) {
        return exec(session -> session.get(resultType, primaryKey));
    }

    public <T> List<T> listAll(
            final String tableName,
            final Class<T> resultType
    ) throws SQLException {
        DBUtils.testName(tableName);
        return exec(
                session -> session
                        .createQuery("FROM " + tableName, resultType)
                        .getResultList()
        ).orElseThrow();
    }

    public <T> void delete(final T toDelete) {
    execConsumer(session -> session.delete(toDelete));
    }

    @Override
    public void close() {
        this.getFactory().close();
    }

    public static void main(final String... args) {
        final Persistence persistence = new Persistence();

        final String sutPK = "username";
        final User sut = new User(sutPK, "password");
        sut.setFullName("fullName");
        sut.setEmail("username@test.edu");
        sut.setShadow(true);
        final String tableName = "cloud_user";

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

        final String certificationID = "certification_id";
        final Certification certification = new Certification();
        certification.setId(certificationID);
        certification.setDescription("Certification Description");
        certification.setPublisher("Certification Publisher");
        certification.setWebsite("website");
        certification.setControls(List.of(control));
        final String certificationTableName = "certification";

        test(
                persistence,
                sut,
                sutPK,
                tableName,
                () -> test(
                        persistence,
                        domain,
                        domainID,
                        domainTableName,
                        () -> test(
                                persistence,
                                control,
                                controlID,
                                controlTableName,
                                () -> test(
                                        persistence,
                                        certification,
                                        certificationID,
                                        certificationTableName,
                                        () -> {}
                                )
                        )
                )
        );
        persistence.close();
    }

    private static <T> void test(
            final Persistence persistence,
            final T sut,
            final Serializable sutPK,
            final String tableName,
            final Runnable runnable
    )  {
        try {
            final Class<T> type = (Class<T>) sut.getClass();
            System.out.println(tableName + persistence.listAll(tableName, type));

            boolean isPresent = persistence.get(type, sutPK).isPresent();
            System.out.println(tableName + " IS_PRESENT: " + isPresent);

            if (!isPresent)
                persistence.save(sut);

            System.out.println(tableName + persistence.listAll(tableName, type));

            var storedValue = persistence.get(type, sutPK);
            System.out.println(tableName + " STORED_VALUE: " + storedValue);

            runnable.run();

            storedValue.ifPresent(persistence::delete);

            isPresent = persistence.get(type, sutPK).isPresent();
            System.out.println(tableName + " IS_PRESENT: " + isPresent);

            System.out.println(tableName + persistence.listAll(tableName, type));
        } catch (SQLException e) {
            e.printStackTrace();
            throw new AssertionError();
        }
    }
}
