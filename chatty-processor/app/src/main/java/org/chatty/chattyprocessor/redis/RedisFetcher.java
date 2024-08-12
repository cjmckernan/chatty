package org.chatty.chattyprocessor.redis;

import java.util.Map;

import redis.clients.jedis.Jedis;

public class RedisFetcher {

  private Jedis jedis;

  public RedisFetcher(String host, int port) {
    jedis = new Jedis(host, port);
  }

  public Map<String, String> fetchUser(String userId) {
    return jedis.hgetAll("user:" + userId);
  }

}
