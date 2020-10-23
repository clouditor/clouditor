package io.clouditor.data_access_layer.utils;

import java.util.function.Consumer;
import java.util.function.Supplier;

/**
 * A functional interface to be able to define lambdas, that could throw a Exception.
 *
 * @param <T> The data type of the supplied value.
 *
 * @author Andreas Hager, andreas.hager@aisec.fraunhofer.de
 */
@FunctionalInterface
public interface BoomingSupplier<T> {

    /**
     * Gets a result.
     *
     * @return a result
     */
    T get() throws Exception;

    default Supplier<T> handleException() {
        return () -> {
            try {
                return this.get();
            } catch (final Exception e) {
                e.printStackTrace();
                throw new AssertionError();
            }
        };
    }

    default Supplier<T> handleException(
            final Consumer<Exception> handleException,
            final Supplier<T> defaultValue
    ) {
        return () -> {
            try {
                return this.get();
            } catch (final Exception exception) {
                handleException.accept(exception);
                return defaultValue.get();
            }
        };
    }
}
