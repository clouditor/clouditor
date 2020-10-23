package io.clouditor.data_access_layer.utils;

import com.google.gson.Gson;
import com.google.gson.JsonSyntaxException;

import java.sql.SQLException;
import java.util.Objects;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * A Utility class for the PostgresDB service.
 *
 * @author Andreas Hager, andreas.hager@aisec.fraunhofer.de
 */
public class DBUtils {

    private static final Gson gson = new Gson();

    private static final String DESCRIPTION_REGEX = "^[\\s0-9a-zA-Z_\\-.]*$";
    private static final Pattern DESCRIPTION_PATTERN = Pattern.compile(DESCRIPTION_REGEX);

    private static final String NAME_REGEX = "^[a-zA-Z][0-9a-zA-Z_\\-]*$";
    private static final Pattern NAME_PATTERN = Pattern.compile(NAME_REGEX);

    private static final String DAY = "(Sun|Mon|Tue|Wed|Thu|Fri|Sat)";
    private static final String MONTH = "(Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec)";
    private static final String DAY_OF_MONTH = "([0-3][0-9]|[0-9])";
    private static final String HOUR = "([0-2][0-9])";
    private static final String MINUTE = "([0-5][0-9])";
    private static final String SECOND = "([0-6][0-9])";
    private static final String TIMEZONE = "([A-Z]*\\s|)";
    private static final String YEAR = "([0-9][0-9][0-9][0-9])";
    private static final String DATE_REGEX = "^"
            + DAY + "\\s"
            + MONTH + "\\s+"
            + DAY_OF_MONTH + "\\s"
            + HOUR + ":"
            + MINUTE + ":"
            + SECOND + "\\s"
            + TIMEZONE
            + YEAR
            + "$";
    private static final Pattern DATE_PATTERN = Pattern.compile(DATE_REGEX);

    private static final String WHERE = " WHERE ";

    public static String getQuery(final SQLMethod sqlMethod, final Table table) {
        Objects.requireNonNull(sqlMethod);
        final String destination;
        if (sqlMethod == SQLMethod.ADD)
            destination = " " + table.getAttributes();
        else // it is a GET or a DELETE Method
            destination = WHERE + table.getKeyQuery();
        return sqlMethod + table.getTableName() + destination + ";";
    }


    public static void testName(final String name) throws SQLException {
        Objects.requireNonNull(name);
        final Matcher nameMatcher = NAME_PATTERN.matcher(name);
        if (!nameMatcher.matches())
            throw new SQLException(
                    "The name: " +
                            name +
                            ", does not match the pattern: " +
                            NAME_REGEX
            );
    }

    public static void testDescription(final String description) throws SQLException {
        Objects.requireNonNull(description);
        final Matcher descriptionMatcher = DESCRIPTION_PATTERN
                .matcher(description);
        if (!descriptionMatcher.matches())
            throw new SQLException(
                    "The description: " + description +
                            ", does not match the pattern: " +
                            DESCRIPTION_REGEX
            );
    }

    public static void testTimeStamp(final String date) throws SQLException {
        Objects.requireNonNull(date);
        final Matcher dateMatcher = DATE_PATTERN.matcher(date);
        if (!dateMatcher.matches())
            throw new SQLException(
                    "The date: " + date +
                            ", does not match the pattern: " +
                            DATE_REGEX
            );
    }

    public static void getJSON(final String json) throws SQLException {
        Objects.requireNonNull(json);
        try {
            gson.fromJson(json, Object.class);
        } catch (JsonSyntaxException exception) {
            throw new SQLException(
                    "The json: " + json +
                            ", is not valid."
            );
        }
    }
}
