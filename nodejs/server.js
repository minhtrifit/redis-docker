const express = require("express");
const axios = require("axios");
const redis = require("redis");
const PORT = 5001;

const app = express();

app.get("/photos", async (req, res) => {
  try {
    const redisClient = redis.createClient();
    await redisClient.connect();

    const value = await redisClient.get("photos");

    if (value) {
      return res.json({
        status: 200,
        message: "Redis cache successfully",
        data: JSON.parse(value),
      });
    } else {
      try {
        const { data } = await axios({
          method: "GET",
          url: "https://jsonplaceholder.typicode.com/todos",
        });

        // Redis cache
        await redisClient.set("photos", JSON.stringify(data));

        return res.json({
          status: 200,
          message: "Get API successfully",
          data: data,
        });
      } catch (error) {
        return res.json({ status: 500, message: "error" });
      }
    }
  } catch (e) {
    return res.status(500);
  }
});

app.listen(PORT, (req, res) => {
  console.log(`Server is runing at port: http://localhost:${PORT}`);
});
