package io.clouditor.data_access_layer.utils;

import java.util.function.Consumer;

/**
 * A functional interface to be able to define lambdas, that could throw a Exception.
 *
 * @author Andreas Hager, andreas.hager@aisec.fraunhofer.de
 */
@FunctionalInterface
public
interface BoomingRunnable {

    /**
     * Run something.
     *
     * @throws Exception if something went wrong
     */
    void run() throws Exception;

    default Runnable handleException() {
        return () -> {
            try {
                this.run();
            } catch (final Exception e) {
                e.printStackTrace();
                throw new AssertionError();
            }
        };
    }

    default Runnable handleException(final Consumer<Exception> handleException) {
        return () -> {
            try {
                this.run();
            } catch (final Exception exception) {
                handleException.accept(exception);
            }
        };
    }
}
