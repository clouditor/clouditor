package io.clouditor.assurance;

import io.clouditor.AbstractEngineUnitTest;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

class CertificationServiceTest extends AbstractEngineUnitTest {

  /* Test Settings */
  @BeforeEach
  @Override
  protected void setUp() {
    super.setUp();

    // init db
    engine.initDB();
    // initialize every else
    engine.init();
  }

  @Override
  protected void cleanUp() {
    super.cleanUp();

    engine.shutdown();
  }

  /* Tests */
  @Test
  void testRemoveCertification() {

    // First assertion: method with id of certification as input
    var certificationService = engine.getService(CertificationService.class);
    Certification mockCert = new Certification();
    mockCert.setId("Mock-Cert-Id");
    certificationService.modifyCertification(mockCert);

    var numberOfCertificationsBeforeRemoval = certificationService.getCertifications().size();
    Assertions.assertNotNull(certificationService.getCertifications().get(mockCert.getId()));

    boolean isRemoved = certificationService.removeCertification(mockCert.getId());
    var numberOfCertificationsAfterRemoval = certificationService.getCertifications().size();

    Assertions.assertTrue(isRemoved);
    Assertions.assertEquals(
        numberOfCertificationsBeforeRemoval - 1, numberOfCertificationsAfterRemoval);
    Assertions.assertNull(certificationService.getCertifications().get(mockCert.getId()));

    // Second assertion: method with certification object as input
    certificationService.modifyCertification(mockCert);
    numberOfCertificationsBeforeRemoval = certificationService.getCertifications().size();
    Assertions.assertNotNull(certificationService.getCertifications().get(mockCert.getId()));

    isRemoved = certificationService.removeCertification(mockCert);
    numberOfCertificationsAfterRemoval = certificationService.getCertifications().size();

    Assertions.assertTrue(isRemoved);
    Assertions.assertEquals(
        numberOfCertificationsBeforeRemoval - 1, numberOfCertificationsAfterRemoval);
    Assertions.assertNull(certificationService.getCertifications().get(mockCert.getId()));

    // Third assertion: method with wrong id of certification as input
    Assertions.assertFalse(certificationService.removeCertification("Wrong_Cert_Id"));
  }

  @Test
  void testRemoveCertifications() {

    // First assertion: method with id of certification as input
    var certificationService = engine.getService(CertificationService.class);
    Certification mockCert1 = new Certification();
    mockCert1.setId("Mock-Cert-Id-1");
    Certification mockCert2 = new Certification();
    mockCert2.setId("Mock-Cert-Id-2");
    certificationService.modifyCertification(mockCert1);
    certificationService.modifyCertification(mockCert2);

    var numberOfCertificationsBeforeRemoval = certificationService.getCertifications().size();
    certificationService.removeAllCertifications();
    var numberOfCertificationsAfterRemoval = certificationService.getCertifications().size();
    Assertions.assertNotEquals(
        numberOfCertificationsBeforeRemoval, numberOfCertificationsAfterRemoval);
    Assertions.assertEquals(0, numberOfCertificationsAfterRemoval);
  }
}
