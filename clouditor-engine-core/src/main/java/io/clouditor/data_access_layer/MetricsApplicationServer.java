package io.clouditor.data_access_layer;

import com.opentable.db.postgres.embedded.EmbeddedPostgres;
import io.clouditor.data_access_layer.utils.BoomingSupplier;
import io.clouditor.data_access_layer.utils.Table;
import io.grpc.Server;
import io.grpc.ServerBuilder;

import java.io.IOException;
import java.util.Arrays;
import java.util.List;
import java.util.Objects;
import java.util.concurrent.TimeUnit;
import java.util.logging.Logger;

public class MetricsApplicationServer {
    // logger for this class is used for the communication in the gRPC connection
    private static final Logger logger = Logger.getLogger(MetricsApplicationServer.class.getName());
    private static final int TIMEOUT_IN_SECONDS = 30;

    // port and the grpc server itself
    private final int port;
    private final Server server;

    public MetricsApplicationServer(final int port, final DBServer metricStorage) {
        this(ServerBuilder.forPort(port), port, metricStorage);
    }

    public MetricsApplicationServer(final ServerBuilder serverBuilder, final int port, final DBServer metricStorage) {
        Objects.requireNonNull(metricStorage);
        if (port <= 0 || port > 65535)
            throw new IllegalArgumentException("The port: " + port + " is not allowed.");
        this.port = port;
        final MetricsApplicationService metricsApplicationService = new MetricsApplicationService(metricStorage);
        this.server = serverBuilder
                .addService(metricsApplicationService)
                .build();
    }

    public void start() throws IOException {
        server.start();
        logger.info("Server started on port " + port);

        Runtime.getRuntime()
                .addShutdownHook(
                        new Thread(() -> {
                            System.err.println("Shut down gRPC server because JVM shuts down");
                            try {
                                MetricsApplicationServer.this.stop();
                            } catch (InterruptedException e) {
                                e.printStackTrace(System.err);
                            }
                            System.err.println("Server shut down");
                        })
                );
    }

    public void stop() throws InterruptedException {
        if (server != null) {
            server.shutdown().awaitTermination(TIMEOUT_IN_SECONDS, TimeUnit.SECONDS);
        }
    }

    public void blockUntilShutdown() throws InterruptedException {
        if (server != null) {
            server.awaitTermination();
        }
    }

    public static void main(final String... args) throws Exception {
        final List<String> arguments = Arrays.asList(args);
        final BoomingSupplier<DBServer> getDBServer = arguments.size() > 1 && "persist".equals(arguments.get(0)) ?
                () -> new PostgresDB("postgres", "postgres", arguments.get(1)) :
                () -> new PostgresDB(
                        Table.init(
                                EmbeddedPostgres.start()
                                        .getPostgresDatabase()
                                        .getConnection()
                        )
                );
        try (
                final DBServer metricStorage = getDBServer
                        .handleException()
                        .get()
        ) {
            final MetricsApplicationServer metricsApplicationServer =
                    new MetricsApplicationServer(50051, metricStorage);
            metricsApplicationServer.start();
            metricsApplicationServer.blockUntilShutdown();
        }
    }
}
