package main

import (
  // "encoding/json"
  "os"
  "fmt"
  "bufio"
  "strings"
  "log"
  "database/sql"

  _ "github.com/go-sql-driver/mysql"
)

var createDB = []string{
  `CREATE DATABASE IF NOT EXISTS F1_GO_MYSQL;`,
  `USE F1_GO_MYSQL;`,
}

var createTable = []string{
  `CREATE TABLE IF NOT EXISTS race_event_directory(
    session_uid BIGINT NOT NULL,
    M_packetFormat YEAR(4),
    packet_version FLOAT(10,6),
    player_car_index TINYINT,
    date DATETIME,
    notes VARCHAR(255),
    PRIMARY KEY (session_uid)
  );`,



  `CREATE TABLE IF NOT EXISTS motion_data (
    id INT NOT NULL AUTO_INCREMENT,
    session_uid BIGINT NOT NULL,
    frame_identifier int NOT NULL,
    suspension_position_rl DECIMAL(15,10),
    suspension_position_rr DECIMAL(15,10),
    suspension_position_fl DECIMAL(15,10),
    suspension_position_fr DECIMAL(15,10),
    M_suspensionVelocity_rl DECIMAL(15,10),
    M_suspensionVelocity_rr DECIMAL(15,10),
    M_suspensionVelocity_fl DECIMAL(15,10),
    M_suspensionVelocity_fr DECIMAL(15,10),
    suspension_acceleration_rl DECIMAL(15,10),
    suspension_acceleration_rr DECIMAL(15,10),
    suspension_acceleration_fl DECIMAL(15,10),
    suspension_acceleration_fr DECIMAL(15,10),
    wheel_speed_rl DECIMAL(15,10),
    wheel_speed_rr DECIMAL(15,10),
    wheel_speed_fl DECIMAL(15,10),
    wheel_speed_fr DECIMAL(15,10),
    wheel_slip_rl DECIMAL(15,10),
    wheel_slip_rr  DECIMAL(15,10),
    wheel_slip_fl DECIMAL(15,10),
    wheel_slip_fr DECIMAL(15,10),
    local_velocity_x DECIMAL(15,10),
    local_velocity_y DECIMAL(15,10),
    local_velocity_z DECIMAL(15,10),
    angular_velocity_x DECIMAL(15,10),
    angular_velocity_y DECIMAL(15,10),
    angular_velocity_z DECIMAL(15,10),
    angular_acceleration_x DECIMAL(15,10),
    angular_acceleration_y DECIMAL(15,10),
    angular_acceleration_z DECIMAL(15,10),
    front_wheels_angle DECIMAL(15,10),
    PRIMARY KEY (id),
    FOREIGN KEY (session_uid) REFERENCES (race_event_directory.session_uid)
  );`,
  `CREATE TABLE IF NOT EXISTS car_motion_data (
    id INT NOT NULL AUTO_INCREMENT,
    motion_packet_id INT NOT NULL,
    m_worldPositionX DECIMAL(15,10),
    m_worldPositionY DECIMAL(15,10),
    m_worldPositionZ DECIMAL(15,10),
    m_worldVelocityX DECIMAL(15,10),
    m_worldVelocityY DECIMAL(15,10),
    m_worldVelocityZ DECIMAL(15,10),
    m_worldForwardDirX SMALLINT,
    m_worldForwardDirY SMALLINT,
    m_worldForwardDirZ SMALLINT,
    m_worldRightDirX SMALLINT,
    m_worldRightDirY SMALLINT,
    m_worldRightDirZ SMALLINT,
    m_gForceLateral DECIMAL(15,10),
    m_gForceLongitudinal DECIMAL(15,10),
    m_gForceVertical DECIMAL(15,10),
    m_yaw DECIMAL(15,10),
    m_pitch DECIMAL(15,10),
    m_roll DECIMAL(15,10),
    PRIMARY KEY (id),
    FOREIGN KEY (motion_packet_id) REFERENCES (motion_packet.id)
  );`,



  `CREATE TABLE IF NOT EXISTS session_data (
    id INT NOT NULL AUTO_INCREMENT,
    session_uid BIGINT NOT NULL,
    frame_identifier int NOT NULL,
    m_weather TINYINT,
    m_trackTemperature TINYINT,
    m_airTemperature TINYINT,
    m_totalLaps TINYINT,
    m_trackLength SMALLINT,
    m_sessionType SMALLINT,
    m_trackId SMALLINT,
    m_era TINYINT,
    m_sessionTimeLeft MEDIUMINT,
    m_sessionDuration MEDIUMINT,
    m_pitSpeedLimit TINYINT,
    m_gamePaused TINYINT,
    m_isSpectating TINYINT,
    m_spectatorCarIndex SMALLINT,
    m_sliProNativeSupport TINYINT,
    m_numMarshalZones TINYINT,
    m_safetyCarStatus TINYINT,
    m_networkGame TINYINT,
    PRIMARY KEY (id),
    FOREIGN KEY (session_uid) REFERENCES (race_event_directory.session_uid)
  );`,
  `CREATE TABLE IF NOT EXISTS marshal_zone (
    id INT NOT NULL AUTO_INCREMENT,
    session_data_id INT NOT NULL,
    m_zoneStart DECIMAL(10,10),
    m_zoneFlag TINYINT,
    PRIMARY KEY (id),
    FOREIGN KEY (session_data_id) REFERENCES (session_data.id)
  );`,



  `CREATE TABLE IF NOT EXISTS lap_data (
    id INT NOT NULL AUTO_INCREMENT,
    session_uid BIGINT NOT NULL,
    frame_identifier int NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (session_uid) REFERENCES (race_event_directory.session_uid)
  );`,
  `CREATE TABLE IF NOT EXISTS car_lap_data (
    id INT NOT NULL AUTO_INCREMENT,
    lap_data_id INT NOT NULL,
    m_lastLapTime DECIMAL(15,10),
    m_currentLapTime DECIMAL(15,10),
    m_bestLapTime DECIMAL(15,10),
    m_sector1Time DECIMAL(15,10),
    m_sector2Time DECIMAL(15,10),
    m_lapDistance DECIMAL(15,10),
    m_totalDistance DECIMAL(10,6),
    m_safetyCarDelta DECIMAL(15,10),
    m_carPosition TINYINT,
    m_currentLapNum TINYINT,
    m_pitStatus TINYINT,
    m_sector TINYINT,
    m_currentLapInvalid TINYINT,
    m_penalties TINYINT,
    m_gridPosition TINYINT,
    m_driverStatus TINYINT,
    m_resultStatus TINYINT,
    PRIMARY KEY (id),
    FOREIGN KEY (lap_data_id) REFERENCES (lap_data.lap_data_id)
  );`,



  `CREATE TABLE IF NOT EXISTS event_data (
    id INT NOT NULL AUTO_INCREMENT,
    session_uid BIGINT NOT NULL,
    frame_identifier int NOT NULL,
    m_eventStringCode CHAR(4),
    PRIMARY KEY (id),
    FOREIGN KEY (session_uid) REFERENCES (race_event_directory.session_uid)
  );`,



  `CREATE TABLE IF NOT EXISTS participant_data (
    id INT NOT NULL AUTO_INCREMENT,
    session_uid BIGINT NOT NULL,
    frame_identifier int NOT NULL,
    m_numCars TINYINT,
    PRIMARY KEY (id),
    FOREIGN KEY (session_uid) REFERENCES (race_event_directory.session_uid)
  );`,
  `CREATE TABLE IF NOT EXISTS car_participant_data (
    id INT NOT NULL AUTO_INCREMENT,
    participant_data_id  INT NOT NULL,
    m_aiControlled TINYINT,
    m_driverId TINYINT,
    m_teamId TINYINT,
    m_raceNumber TINYINT,
    m_nationality TINYINT,
    m_name VARCHAR(48),
    PRIMARY KEY (id),
    FOREIGN KEY (participant_data_id) REFERENCES (participant_data.id)
  );`,



  `CREATE TABLE IF NOT EXISTS setup_data (
    id INT NOT NULL AUTO_INCREMENT,
    session_uid BIGINT NOT NULL,
    frame_identifier int NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (session_uid) REFERENCES (race_event_directory.session_uid)
  );`,
  `CREATE TABLE IF NOT EXISTS car_setup_data (
    id INT NOT NULL AUTO_INCREMENT,
    setup_data_id INT NOT NULL,
    m_frontWing DECIMAL(3,1),
    m_rearWing DECIMAL(3,1),
    m_onThrottle TINYINT,
    m_offThrottle TINYINT,
    m_frontCamber DECIMAL(4,2),
    m_rearCamber DECIMAL(4,2),
    m_frontToe DECIMAL(12,10),
    m_rearToe DECIMAL(12,10),
    m_frontSuspension DECIMAL(3,1),
    m_rearSuspension DECIMAL(3,1),
    m_frontAntiRollBar DECIMAL(3,1),
    m_rearAntiRollBar DECIMAL(3,1),
    m_frontSuspensionHeight DECIMAL(3,1),
    m_rearSuspensionHeight DECIMAL(3,1),
    m_brakePressure TINYINT,
    m_brakeBias TINYINT,
    m_frontTyrePressure DECIMAL(4,2),
    m_rearTyrePressure DECIMAL(4,2),
    m_ballast DECIMAL(3,1),
    m_fuelLoad DECIMAL(3,1),
    PRIMARY KEY (id),
    FOREIGN KEY (setup_data_id) REFERENCES (setup_data.id)
  );`,



  `CREATE TABLE IF NOT EXISTS telemetry_data (
    id INT NOT NULL AUTO_INCREMENT,
    session_uid BIGINT NOT NULL,
    frame_identifier int NOT NULL,
    m_buttonStatus BIT(4),
    PRIMARY KEY (id),
    FOREIGN KEY (session_uid) REFERENCES (race_event_directory.session_uid)
  );`,
  `CREATE TABLE IF NOT EXISTS car_telemetry_data (
    id INT NOT NULL AUTO_INCREMENT,
    telemetry_data_id INT NOT NULL,
    m_speed SMALLINT,
    m_throttle TINYINT,
    m_steer TINYINT,
    m_brake TINYINT,
    m_clutch TINYINT,
    m_gear TINYINT,
    m_engineRPM SMALLINT,
    m_drs TINYINT,
    m_revLightsPercent TINYINT,
    m_brakesTemperature_rl SMALLINT,
    m_brakesTemperature_rr SMALLINT,
    m_brakesTemperature_fl SMALLINT,
    m_brakesTemperature_fr SMALLINT,
    m_tyresSurfaceTemperature_rl SMALLINT,
    m_tyresSurfaceTemperature_rr SMALLINT,
    m_tyresSurfaceTemperature_fl SMALLINT,
    m_tyresSurfaceTemperature_fr SMALLINT,
    m_tyresInnerTemperature_rl SMALLINT,
    m_tyresInnerTemperature_rr SMALLINT,
    m_tyresInnerTemperature_fl SMALLINT,
    m_tyresInnerTemperature_fr SMALLINT,
    m_engineTemperature SMALLINT,
    m_tyresPressure_rl DECIMAL(5,2),
    m_tyresPressure_rr DECIMAL(5,2),
    m_tyresPressure_fl DECIMAL(5,2),
    m_tyresPressure_fr DECIMAL(5,2),
    PRIMARY KEY (id),
    FOREIGN KEY (telemetry_data_id) REFERENCES (telemetry_data.id)
  );`,



  `CREATE TABLE IF NOT EXISTS status_data (
    id INT NOT NULL AUTO_INCREMENT,
    session_uid BIGINT NOT NULL,
    frame_identifier int NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (session_uid) REFERENCES (race_event_directory.session_uid)
  );`,
  `CREATE TABLE IF NOT EXISTS car_status_data (
    id INT NOT NULL AUTO_INCREMENT,
    status_data_id INT NOT NULL,
    m_tractionControl TINYINT,
    m_antiLockBrakes TINYINT,
    m_fuelMix TINYINT,
    m_frontBrakeBias TINYINT,
    m_pitLimiterStatus TINYINT,
    m_fuelInTank DECIMAL(10,7),
    m_fuelCapacity SMALLINT,
    m_maxRPM SMALLINT,
    m_idleRPM SMALLINT,
    m_maxGears TINYINT,
    m_drsAllowed TINYINT,
    m_tyresWear_rl TINYINT,
    m_tyresWear_rr TINYINT,
    m_tyresWear_fl TINYINT,
    m_tyresWear_fr TINYINT,
    m_tyreCompound TINYINT,
    m_tyresDamage_rl TINYINT,
    m_tyresDamage_rr TINYINT,
    m_tyresDamage_fl TINYINT,
    m_tyresDamage_fr TINYINT,
    m_frontLeftWingDamage TINYINT,
    m_frontRightWingDamage TINYINT,
    m_rearWingDamage TINYINT,
    m_engineDamage TINYINT,
    m_gearBoxDamage TINYINT,
    m_exhaustDamage TINYINT,
    m_vehicleFiaFlags TINYINT,
    m_ersStoreEnergy DECIMAL(18,10),
    m_ersDeployMode TINYINT,
    m_ersHarvestedThisLapMGUK DECIMAL(18,10),
    m_ersHarvestedThisLapMGUH DECIMAL(18,10),
    m_ersDeployedThisLap DECIMAL(18,10),
    PRIMARY KEY (id),
    FOREIGN KEY (status_data_id) REFERENCES (status_data.id)
  );`,
}


func deleteDatabase(db *sql.DB) (error) {
   _,err := db.Exec("DROP DATABASE F1_GO_MYSQL")
   if err != nil {
     db.Close()
     return err
   }

   db.Close()
   return nil
}

func createDatabase(db *sql.DB) (*sql.DB, error) {
  for _, stmt := range createDB {
    _, err := db.Exec(stmt)
    if err != nil {
      db.Close()
      return db, err
    }
  }

  // We are now finished making our tables
  // Close the connection and return success!
  return db, nil
}


func createTables(db *sql.DB) (error){
  for _, stmt := range createTable {
    _, err := db.Exec(stmt)
    if err != nil {
      db.Close()
      return err
    }
  }

  // We are now finished making our tables
  // Close the connection and return success!
  db.Close()
  return nil
}

func main() {
  db, err := sql.Open("mysql", "admin:admin@tcp(127.0.0.1:3306)/")
  if err != nil {
     log.Println(err)
  } else {

    arg := os.Args

    if len(arg) > 1 && len(arg) < 3{
      arg = os.Args[1:]  // Since the first postion is the path, our args start at postion 1

      switch(arg[0]) {
      case "create":
        db, err = createDatabase(db)
        if err != nil {
          fmt.Println("Error creating F1_GO database!")
          log.Println(err)
          fmt.Println("Exiting...")
          os.Exit(1)
        }

        fmt.Println("Successfully created F1_GO database")
        fmt.Println("Creating database schema...")

        err = createTables(db)
        if err != nil {
          fmt.Println("Error creating F1_GO schema!")
          log.Println(err)
          fmt.Println("Exiting...")
          os.Exit(1)
        }

        fmt.Println("Successfully created F1_GO schema")
        fmt.Println("F1_GO MYSQL database now ready for use by F1_GO!")
        fmt.Println("Exiting...")

      case "delete":
        reader := bufio.NewReader(os.Stdin)
        fmt.Print("\n!!!WARNING!!!")
        fmt.Print("You are tring to delete your F1_GO MYSQL database")
        fmt.Print("Do you wish to continue? This step can not be reveresed")
        fmt.Print("[ Yes ] or [ No ]: ")
        text, _ := reader.ReadString('\n')

        switch(strings.ToLower(text)) {
        case "no":
          fmt.Println("You decided not to delete the F1_GO MYSQL database")
          fmt.Println("Exiting...")
        case "yes":
          fmt.Println("You have decided to delete the F1_GO MYSQL database")
          fmt.Println("Deleting....")
          err := deleteDatabase(db)

          if err != nil {
            fmt.Println("Error with deleting F1_GO database")
            log.Println(err)
            fmt.Println("Exiting...")
          } else {
            fmt.Println("Successful in deleting F1_GO database")
            fmt.Println("Exiting...")
          }
        default:
          fmt.Println("Incorrect responce! Please try again and answer with a [ yes ] or a [ no ]")
          os.Exit(1)
        }
        fmt.Println(text)

      default:
        log.Println("Incorrect usage! Please use the following formats and try agian")
        fmt.Println("                 1. go run SQLDatabase.go [ create ]")
        fmt.Println("                 2. go run SQLDatabase.go [ delete ]")
      }
    } else {
      log.Println("Incorrect usage! Please use the following formats and try agian")
      fmt.Println("                 1. go run SQLDatabase.go [ create ]")
      fmt.Println("                 2. go run SQLDatabase.go [ delete ]")
    }
  }
}
