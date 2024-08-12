package org.chatty.chattyprocessor.db;

import java.sql.Connection;
import java.sql.ResultSet;
import java.sql.SQLException;
import java.sql.Statement;

/**
 * PostgresRetriever
 */
public class PostgresRetriever {
  private Connection connection;

  public PostgresRetriever(Connection connection) {
    this.connection = connection;
  }

  public String retrieveChats() {
    StringBuilder result = new StringBuilder();
    String query = "SELECT id, user_id, content, created_at FROM message ORDER by created_at DESC";
    try (Statement statement = connection.createStatement(); ResultSet rs = statement.executeQuery(query)) {
      while (rs.next()) {
        int id = rs.getInt("id");
        int userId = rs.getInt("user_id");
        String content = rs.getString("content");
        String createdAt = rs.getTimestamp("created_at").toString();

        result.append("Message ID: ").append(id)
            .append(", User ID: ").append(userId)
            .append(", Content: ").append(content)
            .append(", Created At: ").append(createdAt)
            .append("\n");

      }
    } catch (SQLException e) {
      e.printStackTrace();
    }
    return result.toString();
  }
}
