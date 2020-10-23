package io.clouditor.data_access_layer;

import io.clouditor.metric_api.pb.*;
import io.clouditor.data_access_layer.utils.BoomingRunnable;
import io.clouditor.data_access_layer.utils.BoomingSupplier;
import io.grpc.Status;
import io.grpc.stub.StreamObserver;

import java.util.Objects;
import java.util.Optional;
import java.util.function.Consumer;
import java.util.logging.Level;
import java.util.logging.LogRecord;
import java.util.logging.Logger;


public class MetricsApplicationService extends MetricsAPIServiceGrpc.MetricsAPIServiceImplBase {
    // logger for this class is used for the communication in the gRPC connection
    private static final Logger logger = Logger.getLogger(MetricsApplicationServer.class.getName());

    // Declaring Storage object for accessing and manipulating metrics in memory
    private final DBServer dbServer;

    private final Consumer<Exception> logError = exception -> getLogException().accept(
            exception.toString()
    );

    public MetricsApplicationService(final DBServer dbServer) {
        Objects.requireNonNull(dbServer);
        this.dbServer = dbServer;
    }

    private static Logger getLogger() {
        return logger;
    }

    private static Consumer<String> getLogException() {
        return logException;
    }

    private DBServer getDBServer() {
        return dbServer;
    }

    private void exec(final BoomingRunnable boomingRunnable) {
        boomingRunnable.handleException(logError).run();
    }

    private <T> Optional<T> getOptional(final BoomingSupplier<Optional<T>> boomingSupplier) {
        return boomingSupplier
                .handleException(logError, Optional::empty)
                .get();
    }

    private <T> void handleEmpty(final StreamObserver<T> responseObserver, final String message) {
        Objects.requireNonNull(message);
        responseObserver.onError(
                Status.NOT_FOUND
                        .withDescription(message)
                        .asRuntimeException()
        );
    }

    private static final Consumer<String> logException = message -> getLogger().log(
            new LogRecord(
                    Level.WARNING,
                    message
            )
    );

    @Override
    public void getScale(
            final GetScaleRequest request,
            final StreamObserver<GetScaleResponse> responseObserver
    ) {
        final int lowerBound = request.getLowerBound();
        final int upperBound = request.getUpperBound();
        final Optional<Scale> scale = getOptional(
                () -> getDBServer().loadScale(lowerBound, upperBound)
        );
        if (scale.isPresent()) {
            final GetScaleResponse response = GetScaleResponse
                    .newBuilder()
                    .setScale(scale.get())
                    .build();
            responseObserver.onNext(response);
            responseObserver.onCompleted();
            getLogger().info("Scale sent.");
        } else handleEmpty(
                responseObserver,
                "Method: GET. Scale not found with lowerBound: " + lowerBound
                        + " and upperBound: " + upperBound + "."
        );
    }

    @Override
    public void listScales(
            final ListScaleRequest request,
            final StreamObserver<ListScaleResponse> responseObserver
    ) {
        Objects.requireNonNull(request);
        final ListScaleResponse response = getOptional(
                () -> Optional.of(
                        ListScaleResponse.newBuilder()
                                .addAllScales(getDBServer().loadAllScales())
                                .build()
                )
        ).orElse(ListScaleResponse.getDefaultInstance());
        responseObserver.onNext(response);
        responseObserver.onCompleted();
        getLogger().info("List of scales sent.");
    }

    @Override
    public void addNewScale(
            final AddNewScaleRequest request,
            final StreamObserver<AddNewScaleResponse> responseObserver
    ) {
        final Scale scale = request.getScale();
        final int lowerBound = scale.getLowerBound();
        final int upperBound = scale.getUpperBound();
        final Optional<Scale> storedScale = getOptional(
                () -> getDBServer()
                        .loadScale(lowerBound, upperBound)
        );
        final AddNewScaleResponse response;
        if (storedScale.isEmpty()) {
            exec(() -> getDBServer().storeScale(scale));
            response = AddNewScaleResponse.newBuilder().build();
            getLogger().info("Scale added.");
        } else {
            response = AddNewScaleResponse.getDefaultInstance();
            getLogger()
                .warning(
                        "Method: ADD. Scale already exist. "
                                + "LowerBound: "
                                + lowerBound
                                + " UpperBound: "
                                + upperBound
                                + "."
                );
        }
        responseObserver.onNext(response);
        responseObserver.onCompleted();
    }

    @Override
    public void deleteScale(
            final DeleteScaleRequest request,
            final StreamObserver<DeleteScaleResponse> responseObserver
    ) {
        final int lowerBound = request.getLowerBound();
        final int upperBound = request.getUpperBound();
        final Optional<Scale> storedScale = getOptional(
                () -> getDBServer()
                        .loadScale(lowerBound, upperBound)
        );
        final DeleteScaleResponse response;
        if (storedScale.isPresent()) {
            exec(() -> getDBServer().deleteScale(lowerBound, upperBound));
            response = DeleteScaleResponse.newBuilder().build();
            getLogger().info("Scale deleted.");
        } else{
            response = DeleteScaleResponse.getDefaultInstance();
            getLogger().warning(
                    "Method: DELETE. Scale does not exist. "
                            + "LowerBound: " + lowerBound
                            + " UpperBound: " + upperBound + "."
            );
        }
        responseObserver.onNext(response);
        responseObserver.onCompleted();
    }

    /**
     * getMetric sends the metric with passed metricName to the client
     *
     * @param request
     * @param responseObserver
     */
    @Override
    public void getMetric(
            final GetMetricRequest request,
            final StreamObserver<GetMetricResponse> responseObserver
    ) {
        // Get the id from the client request
        final String metricName = request.getMetricName();

        // Look for metric with passed metricName in the metricStorage
        final Optional<Metric> metric = getOptional(
                () -> getDBServer().loadMetric(metricName)
        );
        // Metric found
        if (metric.isPresent()) {
            // Adding metric to the response builder and build it
            final GetMetricResponse response = GetMetricResponse
                    .newBuilder()
                    .setMetric(metric.get())
                    .build();

            // Send the list, complete the connection and inform the user (server side)
            responseObserver.onNext(response);
            responseObserver.onCompleted();
            getLogger().info("Metric sent");
        } else handleEmpty(
                responseObserver,
                "Method: GET. Metric not found. MetricName: " + metricName + "."
        );
    }


    /**
     * listMetrics sends the list of metrics to the client
     *
     * @param request
     * @param responseObserver
     */
    @Override
    public void listMetrics(
            final ListMetricsRequest request,
            final StreamObserver<ListMetricsResponse> responseObserver
    ) {
        Objects.requireNonNull(request);
        // Adding all metrics of the list of metrics to the response builder and build it
        final ListMetricsResponse response = getOptional(
                () -> Optional.of(
                        ListMetricsResponse.newBuilder()
                                .addAllMetrics(getDBServer().loadAllMetrics())
                                .build()
                )
        ).orElse(ListMetricsResponse.getDefaultInstance());
        // Send the list, complete the connection and inform the user (server side)
        responseObserver.onNext(response);
        responseObserver.onCompleted();
        getLogger().info("List of metrics sent.");
    }


    /**
     * addNewMetric adds the passed metric to the storage
     * @param request
     * @param responseObserver
     */
    @Override
    public void addNewMetric(
            final AddNewMetricRequest request,
            final StreamObserver<AddNewMetricResponse> responseObserver
    ) {
        final Metric metric = request.getMetric();
        final String metricName = metric.getMetricName();
        final int lowerBound = metric.getLowerBound();
        final int upperBound = metric.getUpperBound();
        final Optional<Scale> storedScale = getOptional(
                () -> getDBServer()
                        .loadScale(lowerBound, upperBound)
        );
        final Optional<Metric> storedMetric = getOptional(
                () -> getDBServer()
                        .loadMetric(metricName)
        );
        final AddNewMetricResponse response;
        if (storedScale.isPresent()) {
            if (storedMetric.isEmpty()) {
                exec(() -> getDBServer().storeMetric(metric));
                response = AddNewMetricResponse.newBuilder().build();
                getLogger().info("Metric added.");
            } else {
                response = AddNewMetricResponse.getDefaultInstance();
                getLogger().warning(
                        "Method: ADD. The metric does already exist. "
                                + "MetricName: " + metricName + "."
                );
            }
        } else {
            response = AddNewMetricResponse.getDefaultInstance();
            getLogger().warning(
                    "Method: ADD. For the metric with metricName: " + metricName
                            + ", is no relating scale present. "
                            + "LowerBound: " + lowerBound
                            + "UpperBound: " + upperBound + "."
            );
        }
        responseObserver.onNext(response);
        responseObserver.onCompleted();
    }


    /**
     * deleteMetric removes metric with passed metric_ID from the storage
     *
     * @param request
     * @param responseObserver
     */
    @Override
    public void deleteMetric(
            final DeleteMetricRequest request,
            final StreamObserver<DeleteMetricResponse> responseObserver
    ) {
        final String metricName = request.getMetricName();
        final Optional<Metric> storedMetric = getOptional(
                () -> getDBServer()
                        .loadMetric(metricName)
        );
        final DeleteMetricResponse response;
        if (storedMetric.isPresent()) {
            exec(() -> getDBServer().deleteMetric(metricName));
            response = DeleteMetricResponse.newBuilder().build();
            getLogger().info("Metric deleted.");
        } else {
            response = DeleteMetricResponse.getDefaultInstance();
            getLogger().warning(
                    "Method: DELETE. The metric does already exist. "
                            + "MetricName: " + metricName + "."
            );
        }
        responseObserver.onNext(response);
        responseObserver.onCompleted();
    }

    @Override
    public void getAsset(
            final GetAssetRequest request,
            final StreamObserver<GetAssetResponse> responseObserver
    ) {
        final String assetName = request.getAssetName();
        final Optional<Asset> asset = getOptional(
                () -> getDBServer().loadAsset(assetName)
        );
        if (asset.isPresent()) {
            final GetAssetResponse response = GetAssetResponse
                    .newBuilder()
                    .setAsset(asset.get())
                    .build();
            responseObserver.onNext(response);
            responseObserver.onCompleted();
            getLogger().info("Asset sent.");
        } else handleEmpty(
                responseObserver,
                "Method GET. Asset not found with assetName: "
                            + assetName + "."
        );
    }

    @Override
    public void listAssets(
            final ListAssetsRequest request,
            final StreamObserver<ListAssetsResponse> responseObserver
    ) {
        Objects.requireNonNull(request);
        final ListAssetsResponse response = getOptional(
                () -> Optional.of(
                        ListAssetsResponse.newBuilder()
                                .addAllAssets(getDBServer().loadAllAssets())
                                .build()
                )
        ).orElse(ListAssetsResponse.getDefaultInstance());
        responseObserver.onNext(response);
        responseObserver.onCompleted();
        getLogger().info("List of assets sent.");
    }

    @Override
    public void addNewAsset(
            final AddNewAssetRequest request,
            final StreamObserver<AddNewAssetResponse> responseObserver
    ) {
        final Asset asset = request.getAsset();
        final String assetName = asset.getAssetName();
        final Optional<Asset> storedAsset = getOptional(
                () -> getDBServer()
                        .loadAsset(assetName)
        );
        final AddNewAssetResponse response;
        if (storedAsset.isEmpty()) {
            exec(() -> getDBServer().storeAsset(asset));
            response = AddNewAssetResponse.newBuilder().build();
            getLogger().info("Asset added.");
        } else {
            response = AddNewAssetResponse.getDefaultInstance();
            getLogger().warning(
                    "Method ADD. Asset does already exist. "
                            + "AssetName: " + assetName + "."
            );
        }
        responseObserver.onNext(response);
        responseObserver.onCompleted();
    }

    @Override
    public void deleteAsset(
            final DeleteAssetRequest request,
            final StreamObserver<DeleteAssetResponse> responseObserver
    ) {
        final String assetName = request.getAssetName();
        final Optional<Asset> storedAsset = getOptional(
                () -> getDBServer().loadAsset(assetName)
        );
        final DeleteAssetResponse response;
        if (storedAsset.isPresent()) {
            exec(() -> getDBServer().deleteAsset(assetName));
            response = DeleteAssetResponse.newBuilder().build();
            getLogger().info("Asset deleted.");
        } else {
            response = DeleteAssetResponse.getDefaultInstance();
            getLogger().warning(
                    "Method DELETE. Asset does not exist. "
                            + "AssetName: " + assetName + "."
            );
        }
        responseObserver.onNext(response);
        responseObserver.onCompleted();
    }

    @Override
    public void getEvaluation(
            final GetEvaluationRequest request,
            final StreamObserver<GetEvaluationResponse> responseObserver
    ) {
        final String metricName = request.getMetricName();
        final String timeStamp = request.getTimeStamp();
        final Optional<Evaluation> evaluation = getOptional(
                () -> getDBServer().loadEvaluation(metricName, timeStamp)
        );
        if (evaluation.isPresent()) {
            final GetEvaluationResponse response = GetEvaluationResponse
                    .newBuilder()
                    .setEvaluation(evaluation.get())
                    .build();
            responseObserver.onNext(response);
            responseObserver.onCompleted();
            getLogger().info("Evaluation sent.");
        } else handleEmpty(
                responseObserver,
                "Evaluation not found with metricName: " + metricName
                        + " and timeStamp" + timeStamp + "."
        );
    }

    @Override
    public void listEvaluation(
            final ListEvaluationRequest request,
            final StreamObserver<ListEvaluationResponse> responseObserver) {
        Objects.requireNonNull(request);
        final ListEvaluationResponse response = getOptional(
                () -> Optional.of(
                        ListEvaluationResponse.newBuilder()
                                .addAllEvaluations(getDBServer().loadAllEvaluations())
                                .build()
                )
        ).orElse(ListEvaluationResponse.getDefaultInstance());
        responseObserver.onNext(response);
        responseObserver.onCompleted();
        getLogger().info("List of evaluations sent.");
    }

    @Override
    public void addNewEvaluation(
            final AddNewEvaluationRequest request,
            final StreamObserver<AddNewEvaluationResponse> responseObserver
    ) {
        final Evaluation evaluation = request.getEvaluation();
        final String metricName = evaluation.getMetricName();
        final String timeStamp = evaluation.getTimeStamp();
        final Optional<Metric> storedMetric = getOptional(
                () -> getDBServer().loadMetric(metricName)
        );
        final Optional<Evaluation> storedEvaluation = getOptional(
                () -> getDBServer().loadEvaluation(metricName, timeStamp)
        );
        AddNewEvaluationResponse response = AddNewEvaluationResponse.getDefaultInstance();
        if (storedMetric.isPresent()) {
            if (storedEvaluation.isEmpty()) {
                exec(() -> getDBServer().storeEvaluation(evaluation));
                response = AddNewEvaluationResponse
                        .newBuilder()
                        .build();
                getLogger().info("Evaluation added.");
            } else {
                getLogger().warning(
                        "Method ADD. The Evaluation does already exist. "
                                + "MetricName: " + metricName
                                + ", timeStamp: " + timeStamp + "."
                );
            }
        } else {
            getLogger().warning(
                    "Method ADD. The relating Metric does not exist."
                            + " For the Evaluation: "
                            + "MetricName: " + metricName
                            + ", timeStamp: " + timeStamp + "."
            );
        }
        responseObserver.onNext(response);
        responseObserver.onCompleted();
    }

    @Override
    public void deleteEvaluation(
            final DeleteEvaluationRequest request,
            final StreamObserver<DeleteEvaluationResponse> responseObserver
    ) {
        final String metricName = request.getMetricName();
        final String timeStamp = request.getTimeStamp();
        final Optional<Evaluation> storedEvaluation = getOptional(
                () -> getDBServer().loadEvaluation(metricName, timeStamp)
        );
        final DeleteEvaluationResponse response;
        if (storedEvaluation.isPresent()) {
            exec(() -> getDBServer().deleteEvaluation(metricName, timeStamp));
            response = DeleteEvaluationResponse.newBuilder().build();
            getLogger().info("Evaluation deleted.");
        } else {
            response = DeleteEvaluationResponse.getDefaultInstance();
            getLogger().info(
                    "Method DELETE. The Evaluation does not exist. "
                            + "MetricName: " + metricName
                            + ", timeStamp: " + timeStamp + "."
            );
        }
        responseObserver.onNext(response);
        responseObserver.onCompleted();
    }

    @Override
    public void getEvaluationToAsset(
            final GetEvaluationToAssetRequest request,
            final StreamObserver<GetEvaluationToAssetResponse> responseObserver
    ) {
        final String metricName = request.getMetricName();
        final String asstName = request.getAssetName();
        final String timeStamp = request.getTimeStamp();
        final Optional<EvaluationToAsset> evaluationToAsset = getOptional(
                () -> getDBServer().loadEvaluationToAsset(metricName, asstName, timeStamp)
        );
        if (evaluationToAsset.isPresent()) {
            final GetEvaluationToAssetResponse response = GetEvaluationToAssetResponse
                    .newBuilder()
                    .setEvaluationToAsset(evaluationToAsset.get())
                    .build();
            responseObserver.onNext(response);
            responseObserver.onCompleted();
            getLogger().info("Evaluation sent.");
        } else handleEmpty(
                responseObserver,
                "Evaluation not found with metricName: " + metricName
                        + " and timeStamp" + timeStamp + "."
        );
    }

    @Override
    public void listEvaluationToAsset(
            final ListEvaluationToAssetRequest request,
            final StreamObserver<ListEvaluationToAssetResponse> responseObserver
    ) {
        Objects.requireNonNull(request);
        final ListEvaluationToAssetResponse response = getOptional(
                () -> Optional.of(
                        ListEvaluationToAssetResponse.newBuilder()
                                .addAllEvaluationsToAssets(
                                        getDBServer().loadAllEvaluationsToAssets()
                                ).build()
                )
        ).orElse(ListEvaluationToAssetResponse.getDefaultInstance());
        responseObserver.onNext(response);
        responseObserver.onCompleted();
        getLogger().info("List of evaluationsToAssets sent.");
    }

    @Override
    public void addNewEvaluationToAsset(
            final AddNewEvaluationToAssetRequest request,
            final StreamObserver<AddNewEvaluationToAssetResponse> responseObserver
    ) {
        final EvaluationToAsset evaluationToAsset = request.getEvaluationToAsset();
        final String metricName = evaluationToAsset.getMetricName();
        final String assetName = evaluationToAsset.getAssetName();
        final String timeStamp = evaluationToAsset.getTimeStamp();
        final Optional<Asset> storedAsset = getOptional(
                () -> getDBServer().loadAsset(assetName)
        );
        final Optional<Evaluation> storedEvaluation = getOptional(
                () -> getDBServer().loadEvaluation(metricName, timeStamp)
        );
        final Optional<EvaluationToAsset> storedEvaluationToAsset = getOptional(
                () -> getDBServer().loadEvaluationToAsset(metricName, assetName, timeStamp)
        );
        AddNewEvaluationToAssetResponse response = AddNewEvaluationToAssetResponse.getDefaultInstance();
        if (storedAsset.isPresent()) {
            if (storedEvaluation.isPresent()) {
                if (storedEvaluationToAsset.isEmpty()) {
                    exec(() -> getDBServer().storeEvaluationToAsset(evaluationToAsset));
                    response = AddNewEvaluationToAssetResponse
                            .newBuilder()
                            .build();
                    getLogger().info("EvaluationToAsset added.");
                } else {
                    getLogger().warning(
                            "Method ADD. The EvaluationToAsset does already exist. "
                                    + "MetricName: " + metricName
                                    + ", AssetName: " + assetName
                                    + ", timeStamp: " + timeStamp + "."
                    );
                }
            } else {
                getLogger().warning(
                        "Method ADD. The Evaluation does already exist. "
                                + "MetricName: " + metricName
                                + ", timeStamp: " + timeStamp + "."
                );
            }
        } else {
            getLogger().warning(
                    "Method ADD. The relating Asset does not exist."
                            + " For the Evaluation: "
                            + "assetName: " + assetName + "."
            );
        }
        responseObserver.onNext(response);
        responseObserver.onCompleted();
    }

    @Override
    public void deleteEvaluationToAsset(
            final DeleteEvaluationToAssetRequest request,
            final StreamObserver<DeleteEvaluationToAssetResponse> responseObserver
    ) {
        final String metricName = request.getMetricName();
        final String assetName = request.getAssetName();
        final String timeStamp = request.getTimeStamp();
        final Optional<EvaluationToAsset> storedEvaluationToAsset = getOptional(
                () -> getDBServer().loadEvaluationToAsset(
                        metricName,
                        assetName,
                        timeStamp
                )
        );
        final DeleteEvaluationToAssetResponse response;
        if (storedEvaluationToAsset.isPresent()) {
            exec(() -> getDBServer().deleteEvaluationToAsset(metricName, assetName, timeStamp));
            response = DeleteEvaluationToAssetResponse.newBuilder().build();
            getLogger().info("EvaluationToAsset deleted.");
        } else {
            response = DeleteEvaluationToAssetResponse.getDefaultInstance();
            getLogger().info(
                    "Method DELETE. The EvaluationToAsset does not exist. "
                            + "MetricName: " + metricName
                            + ", assetName: " + assetName
                            + ", timeStamp: " + timeStamp + "."
            );
        }
        responseObserver.onNext(response);
        responseObserver.onCompleted();
    }
}