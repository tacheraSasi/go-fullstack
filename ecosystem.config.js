module.exports = {
  apps: [{
    name: "go-fullstack",
    script: "./bin/api",
    instances: 1,
    autostart: true,
    watch: false,
    max_restarts: 10,
    restart_delay: 2000,
    env: {
      GIN_MODE: "release",
      SERVER_PORT: "8080",
      DB_TYPE: "sqlite",
      DB_PATH: "./data/core.db",
      JWT_SECRET: "change-this-to-a-random-secret",
      JWT_EXPIRES_IN: "24",
      LOG_FILE_PATH: "./logs/app.log",
    },
    error_file: "./logs/pm2-error.log",
    out_file: "./logs/pm2-out.log",
    merge_logs: true,
    log_date_format: "YYYY-MM-DD HH:mm:ss Z",
  }]
};
