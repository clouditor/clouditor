package io.clouditor.data_access_layer;

import io.clouditor.AbstractEngineUnitTest;
import org.junit.jupiter.api.BeforeEach;

public class PersistenceManagerTest extends AbstractEngineUnitTest {

    @Override
    @BeforeEach
    protected void setUp() {
        super.setUp();

        this.engine.initDB();
    }



}
