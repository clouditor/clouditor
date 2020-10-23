package io.clouditor.data_access_layer.dbserver.utils;

import io.clouditor.data_access_layer.utils.BoomingFunction;
import io.clouditor.data_access_layer.utils.DBUtils;
import io.clouditor.data_access_layer.utils.SQLMethod;
import io.clouditor.data_access_layer.utils.Table;
import org.junit.jupiter.api.Test;

import java.sql.SQLException;
import java.util.Date;

import static org.junit.Assert.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;

public class DBUtilsTest {

    @Test
    public void testGetQueryWithNullMethod() {
        // arrange // act
        assertThrows(
                NullPointerException.class,
                () -> DBUtils.getQuery(null, Table.ASSET)
        );
    }

    @Test
    public void testGetQueryWithNullTable() {
        // arrange // act
        assertThrows(
                NullPointerException.class,
                () -> DBUtils.getQuery(SQLMethod.GET, null)
        );
    }

    @Test
    public void getQueryForDELETEScale() {
        // arrange
        final SQLMethod sqlMethod = SQLMethod.DELETE;
        final Table table = Table.SCALE;
        final String want = "DELETE FROM scale WHERE lower_bound=? AND upper_bound=?;";
        // act
        final String have = DBUtils.getQuery(sqlMethod, table);
        // assert
        assertEquals(want, have);
    }

    @Test
    public void getQueryForDELETEAsset() {
        // arrange
        final SQLMethod sqlMethod = SQLMethod.DELETE;
        final Table table = Table.ASSET;
        final String want = "DELETE FROM asset WHERE asset_name=?;";
        // act
        final String have = DBUtils.getQuery(sqlMethod, table);
        // assert
        assertEquals(want, have);
    }

    @Test
    public void getQueryForDeleteMetric() {
        // arrange
        final SQLMethod sqlMethod = SQLMethod.DELETE;
        final Table table = Table.METRIC;
        final String want = "DELETE FROM metric WHERE metric_name=?;";
        // act
        final String have = DBUtils.getQuery(sqlMethod, table);
        // assert
        assertEquals(want, have);
    }

    @Test
    public void getQueryForDeleteEvaluation() {
        // arrange
        final SQLMethod sqlMethod = SQLMethod.DELETE;
        final Table table = Table.EVALUATION;
        final String want = "DELETE FROM evaluation WHERE metric_name=? AND time_stamp=?;";
        // act
        final String have = DBUtils.getQuery(sqlMethod, table);
        // assert
        assertEquals(want, have);
    }

    @Test
    public void getQueryForAddScale() {
        // arrange
        final SQLMethod sqlMethod = SQLMethod.ADD;
        final Table table = Table.SCALE;
        final String want = "INSERT INTO scale values(?,?,?);";
        // act
        final String have = DBUtils.getQuery(sqlMethod, table);
        // assert
        assertEquals(want, have);
    }

    @Test
    public void getQueryForAddAsset() {
        // arrange
        final SQLMethod sqlMethod = SQLMethod.ADD;
        final Table table = Table.ASSET;
        final String want = "INSERT INTO asset values(?,?);";
        // act
        final String have = DBUtils.getQuery(sqlMethod, table);
        // assert
        assertEquals(want, have);
    }

    @Test
    public void getQueryForAddMetric() {
        // arrange
        final SQLMethod sqlMethod = SQLMethod.ADD;
        final Table table = Table.METRIC;
        final String want = "INSERT INTO metric values(?,?,?,?);";
        // act
        final String have = DBUtils.getQuery(sqlMethod, table);
        // assert
        assertEquals(want, have);
    }

    @Test
    public void getQueryForAddEvaluation() {
        // arrange
        final SQLMethod sqlMethod = SQLMethod.ADD;
        final Table table = Table.EVALUATION;
        final String want = "INSERT INTO evaluation values(?,?,?);";
        // act
        final String have = DBUtils.getQuery(sqlMethod, table);
        // assert
        assertEquals(want, have);
    }

    @Test
    public void testNameNull()  {
        // arrange // act
        assertThrows(
                NullPointerException.class,
                () -> DBUtils.testName(null)
        );
    }

    @Test
    public void testNameEmptyString()  {
        // arrange
        final String sut = "";
        // act
        assertThrows(
                SQLException.class,
                () -> DBUtils.testName(sut)
        );
    }

    @Test
    public void testNameOnlyWhitespace()  {
        // arrange
        final String sut = " ";
        // act
        assertThrows(
                SQLException.class,
                () -> DBUtils.testName(sut)
        );
    }

    @Test
    public void testNameStringContainingWhitespace()  {
        // arrange
        final String sut = "bla bulb";
        // act
        assertThrows(
                SQLException.class,
                () -> DBUtils.testName(sut)
        );
    }

    @Test
    public void testNameStringContainingNewLine()  {
        // arrange
        final String sut = "bla\r\nbulb";
        // act
        assertThrows(
                SQLException.class,
                () -> DBUtils.testName(sut)
        );
    }

    @Test
    public void testNameStringContainingTab()  {
        // arrange
        final String sut = "bla\tbulb";
        // act
        assertThrows(
                SQLException.class,
                () -> DBUtils.testName(sut)
        );
    }

    @Test
    public void testNameLeadingNumber()  {
        // arrange
        final String sut = "1bla_bulb";
        // act
        assertThrows(SQLException.class, () -> DBUtils.testName(sut));
    }

    @Test
    public void testNameLeadingUnderscore()  {
        // arrange
        final String sut = "_bla_bulb";
        // act
        assertThrows(SQLException.class, () -> DBUtils.testName(sut));
    }


    @Test
    public void testNameLeadingScore()  {
        // arrange
        final String sut = "-bla_bulb";
        // act
        assertThrows(SQLException.class, () -> DBUtils.testName(sut));
    }

    @Test
    public void testNameContainingSemicolon()  {
        // arrange
        final String sut = "bla_bulb;";
        // act
        assertThrows(SQLException.class, () -> DBUtils.testName(sut));
    }

    @Test
    public void testNameContainingApostrophe()  {
        // arrange
        final String sut = "bla_bulb'";
        // act
        assertThrows(SQLException.class, () -> DBUtils.testName(sut));
    }

    @Test
    public void testNameContainingNumbersUnderscoresScoresETC() throws SQLException {
        // arrange
        final String sut = "bla_bulb-18941_asdf_qwer";
        // act
        DBUtils.testName(sut);
    }

    @Test
    public void testNameContainingNumbersUnderscoresScoresETC2() throws SQLException {
        // arrange
        final String sut = "Bbla_bul346SF--__5678756_Gb-18941_asaAFdf_qwSGDer";
        // act
        DBUtils.testName(sut);
    }

    @Test
    public void testDescriptionNull()  {
        // arrange // act
        assertThrows(NullPointerException.class, () -> DBUtils.testDescription(null));
    }

    @Test
    public void testDescriptionSemicolon()  {
        // arrange
        final String sut = ";";
        // act
        assertThrows(SQLException.class, () -> DBUtils.testDescription(sut));
    }

    @Test
    public void testDescriptionApostrophe()  {
        // arrange
        final String sut = "'";
        // act
        assertThrows(SQLException.class, () -> DBUtils.testDescription(sut));
    }

    @Test
    public void testDescriptionWhitespace() throws SQLException {
        // arrange
        final String sut = " ";
        // act
        DBUtils.testDescription(sut);
    }

    @Test
    public void testDescriptionNumber() throws SQLException {
        // arrange
        final String sut = "5";
        // act
        DBUtils.testDescription(sut);
    }

    @Test
    public void testDescriptionNewLine() throws SQLException {
        // arrange
        final String sut = "\r\n";
        // act
        DBUtils.testDescription(sut);
    }

    @Test
    public void testDescriptionTabulator() throws SQLException {
        // arrange
        final String sut = "\t";
        // act
        DBUtils.testDescription(sut);
    }

    @Test
    public void testDescriptionAll() throws SQLException {
        // arrange
        final String sut = "2346 JHV _oiag5678............\t qwer\r\n ----- ---- -___  asdf\t";
        // act
        DBUtils.testDescription(sut);
    }

    @Test
    public void testTimeStampNullBooms()  {
        //arrange //act
        assertThrows(NullPointerException.class, () -> DBUtils.testTimeStamp(null));
    }

    @Test
    public void testTimeStampSomeStringBooms()  {
        //arrange
        final String sut = "asdf qwer ";
        // act
        assertThrows(SQLException.class, () -> DBUtils.testTimeStamp(sut));
    }

    @Test
    public void testTimeStampWithDateExplicitExampleMinBooms()  {
        //arrange
        final String sut = "Mon Jan 00 00:00:00 AF0000";
        // act
        assertThrows(SQLException.class, () -> DBUtils.testTimeStamp(sut));
    }

    @Test
    public void testTimeStampWithDateExplicitExampleMinBooms2()  {
        //arrange
        final String sut = "Mon Jan 0000:00:00 AF 0000";
        // act
        assertThrows(SQLException.class, () -> DBUtils.testTimeStamp(sut));
    }

    @Test
    public void testTimeStampWithDateExplicitExampleBooms()  {
        //arrange
        final String sut = "Mon Jan 00 00:00:00 AF 00000";
        // act
        assertThrows(SQLException.class, () -> DBUtils.testTimeStamp(sut));
    }

    @Test
    public void testTimeStampWithDateExplicitExampleBooms2()  {
        //arrange
        final String sut = "Mon Jan 00 00:00:00 AF 0000 ";
        // act
        assertThrows(SQLException.class, () -> DBUtils.testTimeStamp(sut));
    }

    @Test
    public void testTimeStampWithDateExplicitExampleBooms3()  {
        //arrange
        final String sut = "\nMon Jan 00 00:00:00 AF 0000";
        // act
        assertThrows(SQLException.class, () -> DBUtils.testTimeStamp(sut));
    }

    @Test
    public void testTimeStampWithDateExplicitExampleBooms4()  {
        //arrange
        final String sut = "\r\nMon Jan 00 00:00:00 AF 0000";
        // act
        assertThrows(SQLException.class, () -> DBUtils.testTimeStamp(sut));
    }

    @Test
    public void testTimeStampWithDateExplicitExampleBooms5()  {
        //arrange
        final String sut = "\tMon Jan 00 00:00:00 AF 0000";
        // act
        assertThrows(SQLException.class, () -> DBUtils.testTimeStamp(sut));
    }

    @Test
    public void testTimeStampWithDateExplicitExampleBooms6()  {
        //arrange
        final String sut = " Mon Jan 00 00:00:00 AF 0000";
        // act
        assertThrows(SQLException.class, () -> DBUtils.testTimeStamp(sut));
    }

    @Test
    public void testTimeStampWithDateExplicitExampleBooms7()  {
        //arrange
        final String sut = "Row Jan 00 00:00:00 AF 0000";
        // act
        assertThrows(SQLException.class, () -> DBUtils.testTimeStamp(sut));
    }

    @Test
    public void testTimeStampWithDateExplicitExampleBooms8()  {
        //arrange
        final String sut = "Mon Fri 00 00:00:00 AF 0000";
        // act
        assertThrows(SQLException.class, () -> DBUtils.testTimeStamp(sut));
    }

    @Test
    public void testTimeStampWithDateExplicitExampleBooms9()  {
        //arrange
        final String sut = "Mon Jan 99 00:00:00 AF 0000";
        // act
        assertThrows(SQLException.class, () -> DBUtils.testTimeStamp(sut));
    }

    @Test
    public void testTimeStampWithDateExplicitExampleBooms10()  {
        //arrange
        final String sut = "Mon Jan 00 99:00:00 AF 0000";
        // act
        assertThrows(SQLException.class, () -> DBUtils.testTimeStamp(sut));
    }

    @Test
    public void testTimeStampWithDateExplicitExampleBooms11()  {
        //arrange
        final String sut = "Mon Jan 00 00:99:00 AF 0000";
        // act
        assertThrows(SQLException.class, () -> DBUtils.testTimeStamp(sut));
    }

    @Test
    public void testTimeStampWithDateExplicitExampleBooms12()  {
        //arrange
        final String sut = "Mon Jan 00 00:00:99 AF 0000";
        // act
        assertThrows(SQLException.class, () -> DBUtils.testTimeStamp(sut));
    }

    @Test
    public void testTimeStampWithDateExplicitExampleBooms13()  {
        //arrange
        final String sut = "Mon Jan 00 00:00:00 af 0000";
        // act
        assertThrows(SQLException.class, () -> DBUtils.testTimeStamp(sut));
    }

    @Test
    public void testTimeStampWithDate() throws SQLException {
        //arrange
        final String sut = new Date().toString();
        System.out.println(sut);
        // act
        DBUtils.testTimeStamp(sut);
    }

    @Test
    public void testTimeStampWithDateExplicitExample() throws SQLException {
        //arrange
        final String sut = "Thu Sep 22 12:05:31 CEST 2020";
        // act
        DBUtils.testTimeStamp(sut);
    }

    @Test
    public void testTimeStampWithDateExplicitExampleMax() throws SQLException {
        //arrange
        final String sut = "Sun Dec 39 29:59:69 CESTDGIUIOGUIFZUKGUFVUFIU 9999";
        // act
        DBUtils.testTimeStamp(sut);
    }

    @Test
    public void testTimeStampWithDateExplicitExampleMin() throws SQLException {
        //arrange
        final String sut = "Mon Jan 00 00:00:00  0000";
        // act
        DBUtils.testTimeStamp(sut);
    }

    @Test
    public void testTimeStampWithDateExplicitExampleMin2() throws SQLException {
        //arrange
        final String sut = "Mon Jan 00 00:00:00 0000";
        // act
        DBUtils.testTimeStamp(sut);
    }

    @Test
    public void testTimeStampWithDateExplicitExampleMin3() throws SQLException {
        //arrange
        final String sut = "Thu Oct  1 16:38:55 2020";
        // act
        DBUtils.testTimeStamp(sut);
    }

    @Test
    public void testBoomingFunction() {
        // arrange
        final String test = "bla";
        final BoomingFunction<String, String> boomingFunction = string -> {
            throw new SQLException("Test exception. EXPECTED!!!");
        };
        // act
        assertThrows(
                AssertionError.class,
                () -> boomingFunction.handleSQLException().apply(test)
        );
    }

    @Test
    public void testNullJSONBooms()  {
        // arrange // act
        assertThrows(NullPointerException.class, () -> DBUtils.getJSON(null));
    }

    @Test
    public void testJSONBooms()  {
        // arrange
        final String sut = "{";
        // act
        assertThrows(SQLException.class, () -> DBUtils.getJSON(sut));
    }

    @Test
    public void testJSONBooms2()  {
        // arrange
        final String sut = ";";
        // act
        assertThrows(SQLException.class, () -> DBUtils.getJSON(sut));
    }

    @Test
    public void testEmptyJSON() throws SQLException {
        // arrange
        final String sut = "";
        // act
        DBUtils.getJSON(sut);
    }

    @Test
    public void testJSON() throws SQLException {
    // arrange
    final String sut =
        "{ \"developers\": [ { \"firstName\":\"John\", \"lastName\":\"von Neumann\" } ] }";
        // act
        DBUtils.getJSON(sut);
    }

}
