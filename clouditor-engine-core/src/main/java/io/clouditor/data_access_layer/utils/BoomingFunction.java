package io.clouditor.data_access_layer.utils;

import java.sql.SQLException;
import java.util.function.Function;

/**
 * A functional interface to be able to define lambdas, that could throw a Exception.
 *
 * @param <T> The domain data type of the function.
 * @param <U> The codomain data type of the function.
 *
 * @author Andreas Hager, andreas.hager@aisec.fraunhofer.de
 */
@FunctionalInterface
public interface BoomingFunction<T, U> {

    /**
     * The single abstract function of the Interface.
     *
     * @param input The input value of the function. It has to be effectively final.
     * @return The output of the function.
     * @throws SQLException If something went wrong.
     */
    U apply(final T input) throws SQLException;

    default Function<T, U> handleSQLException() {
        return input -> {
            try{
                return apply(input);
            } catch (final SQLException sqlException) {
                sqlException.printStackTrace();
                throw new AssertionError();
            }
        };
    }
}
