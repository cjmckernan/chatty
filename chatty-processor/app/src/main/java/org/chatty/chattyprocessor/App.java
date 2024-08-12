/*
 * Redis processor for data storage
 */
package org.chatty.chattyprocessor;

import org.chatty.chattyprocessor.db.PostgresConnector;

public class App {

  public PostgresConnector conn;

  public String getGreeting() {
    return "Hello, welcome to Chatty Processor!";
  }

  public static void main(String[] args) {

  }
}
