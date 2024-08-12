package org.chatty.chattyprocessor.db;

import java.sql.Connection;
import java.sql.DriverManager;
import java.sql.SQLException;

/**
 * PostgresConnector
 */
public class PostgresConnector {

  private static final String URL = "jdbc:posgresql://localhost:5432/chatty";
  private static final String USER = "posgres";
  private static final String PASSWORD = "postgrespassword";

  public Connection connect() throws SQLException {
    return DriverManager.getConnection(URL, USER, PASSWORD);
  }
}
