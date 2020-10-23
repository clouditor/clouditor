package io.clouditor.data_access_layer.utils;

/**
 * A enumeration of all used methods in the PostgresDB service.
 *
 * @author Andreas Hager, andreas.hager@aisec.fraunhofer.de
 */
public enum SQLMethod {
    GET("SELECT * FROM "),
    DELETE("DELETE FROM "),
    ADD("INSERT INTO ");

    private final String query;

    SQLMethod(final String query) {
        this.query = query;
    }

    @Override
    public String toString() {
        return query;
    }
}
