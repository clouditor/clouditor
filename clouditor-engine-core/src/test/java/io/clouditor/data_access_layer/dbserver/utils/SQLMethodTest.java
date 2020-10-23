package io.clouditor.data_access_layer.dbserver.utils;

import io.clouditor.data_access_layer.utils.SQLMethod;
import org.junit.jupiter.api.Test;

import static org.junit.Assert.assertEquals;

public class SQLMethodTest {

    @Test
    public void testGET() {
        // arrange
        final String want = "SELECT * FROM ";
        // act
        final String have = SQLMethod.GET.toString();
        // assert
        assertEquals(want, have);
    }

    @Test
    public void testADD() {
        // arrange
        final String want = "INSERT INTO ";
        // act
        final String have = SQLMethod.ADD.toString();
        // assert
        assertEquals(want, have);
    }

    @Test
    public void testDELETE() {
        // arrange
        final String want = "DELETE FROM ";
        // act
        final String have = SQLMethod.DELETE.toString();
        // assert
        assertEquals(want, have);
    }
}
