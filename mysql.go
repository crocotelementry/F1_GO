package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/crocotelementry/F1_GO/structs"
	"github.com/fatih/color"
	"github.com/go-sql-driver/mysql"
)

var (
	saved_mysql_password     = ""
	mysql_login_string_front = "root:"
	mysql_login_string_back  = "@tcp(127.0.0.1:3306)/"
)

var createDB = []string{
	`CREATE DATABASE IF NOT EXISTS F1_GO_MYSQL;`,
	`USE F1_GO_MYSQL;`,
}

var tableNames = []string{
	`race_event_directory`,
	`motion_data`,
	`car_motion_data`,
	`session_data`,
	`marshal_zone`,
	`lap_data`,
	`car_lap_data`,
	`event_data`,
	`participant_data`,
	`car_participant_data`,
	`setup_data`,
	`car_setup_data`,
	`telemetry_data`,
	`car_telemetry_data`,
	`status_data`,
	`car_status_data`,
}

var createTables = []string{
	`                               CREATE TABLE IF NOT EXISTS race_event_directory(
                                   session_uid BIGINT UNSIGNED NOT NULL,
                                   M_packetFormat YEAR(4),
                                   packet_version FLOAT(10,6),
                                   player_car_index TINYINT,
																	 session_start DATETIME,
																	 session_end DATETIME,
                                   notes VARCHAR(255),
                                   PRIMARY KEY (session_uid)
                                 );


                                 `,

	`                               CREATE TABLE IF NOT EXISTS motion_data (
                                   id INT NOT NULL AUTO_INCREMENT,
                                   session_uid BIGINT UNSIGNED NOT NULL,
                                   frame_identifier INT NOT NULL,
																	 session_time DECIMAL(16,10),
                                   suspension_position_rl DECIMAL(16,10),
                                   suspension_position_rr DECIMAL(16,10),
                                   suspension_position_fl DECIMAL(16,10),
                                   suspension_position_fr DECIMAL(16,10),
                                   M_suspensionVelocity_rl DECIMAL(16,10),
                                   M_suspensionVelocity_rr DECIMAL(16,10),
                                   M_suspensionVelocity_fl DECIMAL(16,10),
                                   M_suspensionVelocity_fr DECIMAL(16,10),
                                   suspension_acceleration_rl DECIMAL(16,10),
                                   suspension_acceleration_rr DECIMAL(16,10),
                                   suspension_acceleration_fl DECIMAL(16,10),
                                   suspension_acceleration_fr DECIMAL(16,10),
                                   wheel_speed_rl DECIMAL(16,10),
                                   wheel_speed_rr DECIMAL(16,10),
                                   wheel_speed_fl DECIMAL(16,10),
                                   wheel_speed_fr DECIMAL(16,10),
                                   wheel_slip_rl DECIMAL(16,10),
                                   wheel_slip_rr  DECIMAL(16,10),
                                   wheel_slip_fl DECIMAL(16,10),
                                   wheel_slip_fr DECIMAL(16,10),
                                   local_velocity_x DECIMAL(16,10),
                                   local_velocity_y DECIMAL(16,10),
                                   local_velocity_z DECIMAL(16,10),
                                   angular_velocity_x DECIMAL(16,10),
                                   angular_velocity_y DECIMAL(16,10),
                                   angular_velocity_z DECIMAL(16,10),
                                   angular_acceleration_x DECIMAL(16,10),
                                   angular_acceleration_y DECIMAL(16,10),
                                   angular_acceleration_z DECIMAL(16,10),
                                   front_wheels_angle DECIMAL(16,10),
                                   PRIMARY KEY (id),
                                   FOREIGN KEY (session_uid) REFERENCES race_event_directory(session_uid)
                                 );`,
	`                               CREATE TABLE IF NOT EXISTS car_motion_data (
                                   id INT NOT NULL AUTO_INCREMENT,
                                   motion_packet_id INT NOT NULL,
																	 car_index				INT NOT NULL,
                                   m_worldPositionX DECIMAL(16,10),
                                   m_worldPositionY DECIMAL(16,10),
                                   m_worldPositionZ DECIMAL(16,10),
                                   m_worldVelocityX DECIMAL(16,10),
                                   m_worldVelocityY DECIMAL(16,10),
                                   m_worldVelocityZ DECIMAL(16,10),
                                   m_worldForwardDirX SMALLINT,
                                   m_worldForwardDirY SMALLINT,
                                   m_worldForwardDirZ SMALLINT,
                                   m_worldRightDirX SMALLINT,
                                   m_worldRightDirY SMALLINT,
                                   m_worldRightDirZ SMALLINT,
                                   m_gForceLateral DECIMAL(16,10),
                                   m_gForceLongitudinal DECIMAL(16,10),
                                   m_gForceVertical DECIMAL(16,10),
                                   m_yaw DECIMAL(16,10),
                                   m_pitch DECIMAL(16,10),
                                   m_roll DECIMAL(16,10),
                                   PRIMARY KEY (id),
                                   FOREIGN KEY (motion_packet_id) REFERENCES motion_data(id)
                                 );


                                 `,
	`                               CREATE TABLE IF NOT EXISTS session_data (
                                   id INT NOT NULL AUTO_INCREMENT,
                                   session_uid BIGINT UNSIGNED NOT NULL,
                                   frame_identifier int NOT NULL,
																	 session_time DECIMAL(16,10),
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
                                   FOREIGN KEY (session_uid) REFERENCES race_event_directory(session_uid)
                                 );`,
	`                               CREATE TABLE IF NOT EXISTS marshal_zone (
                                   id INT NOT NULL AUTO_INCREMENT,
                                   session_data_id INT NOT NULL,
																	 car_index INT NOT NULL,
                                   m_zoneStart DECIMAL(10,10),
                                   m_zoneFlag TINYINT,
                                   PRIMARY KEY (id),
                                   FOREIGN KEY (session_data_id) REFERENCES session_data(id)
                                 );


                                 `,
	`                               CREATE TABLE IF NOT EXISTS lap_data (
                                   id INT NOT NULL AUTO_INCREMENT,
                                   session_uid BIGINT UNSIGNED NOT NULL,
                                   frame_identifier int NOT NULL,
																	 session_time DECIMAL(16,10),
                                   PRIMARY KEY (id),
                                   FOREIGN KEY (session_uid) REFERENCES race_event_directory(session_uid)
                                 );`,
	`                               CREATE TABLE IF NOT EXISTS car_lap_data (
                                   id INT NOT NULL AUTO_INCREMENT,
                                   lap_data_id INT NOT NULL,
																	 car_index INT NOT NULL,
                                   m_lastLapTime DECIMAL(16,10),
                                   m_currentLapTime DECIMAL(16,10),
                                   m_bestLapTime DECIMAL(16,10),
                                   m_sector1Time DECIMAL(16,10),
                                   m_sector2Time DECIMAL(16,10),
                                   m_lapDistance DECIMAL(16,10),
                                   m_totalDistance DECIMAL(16,10),
                                   m_safetyCarDelta DECIMAL(16,10),
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
                                   FOREIGN KEY (lap_data_id) REFERENCES lap_data(id)
                                 );


                                 `,
	`                               CREATE TABLE IF NOT EXISTS event_data (
                                   id INT NOT NULL AUTO_INCREMENT,
                                   session_uid BIGINT UNSIGNED NOT NULL,
                                   frame_identifier int NOT NULL,
																	 session_time DECIMAL(16,10),
                                   m_eventStringCode CHAR(4),
                                   PRIMARY KEY (id),
                                   FOREIGN KEY (session_uid) REFERENCES race_event_directory(session_uid)
                                 );


                                 `,
	`                               CREATE TABLE IF NOT EXISTS participant_data (
                                   id INT NOT NULL AUTO_INCREMENT,
                                   session_uid BIGINT UNSIGNED NOT NULL,
                                   frame_identifier int NOT NULL,
																	 session_time DECIMAL(16,10),
                                   m_numCars TINYINT,
                                   PRIMARY KEY (id),
                                   FOREIGN KEY (session_uid) REFERENCES race_event_directory(session_uid)
                                 );`,
	`                               CREATE TABLE IF NOT EXISTS car_participant_data (
                                   id INT NOT NULL AUTO_INCREMENT,
                                   participant_data_id  INT NOT NULL,
																	 car_index INT NOT NULL,
                                   m_aiControlled TINYINT,
                                   m_driverId TINYINT,
                                   m_teamId TINYINT,
                                   m_raceNumber TINYINT,
                                   m_nationality TINYINT,
                                   m_name VARCHAR(48),
                                   PRIMARY KEY (id),
                                   FOREIGN KEY (participant_data_id) REFERENCES participant_data(id)
                                 );


                                 `,
	`                               CREATE TABLE IF NOT EXISTS setup_data (
                                   id INT NOT NULL AUTO_INCREMENT,
                                   session_uid BIGINT UNSIGNED NOT NULL,
                                   frame_identifier int NOT NULL,
																	 session_time DECIMAL(16,10),
                                   PRIMARY KEY (id),
                                   FOREIGN KEY (session_uid) REFERENCES race_event_directory(session_uid)
                                 );`,
	`                               CREATE TABLE IF NOT EXISTS car_setup_data (
                                   id INT NOT NULL AUTO_INCREMENT,
                                   setup_data_id INT NOT NULL,
																	 car_index INT NOT NULL,
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
                                   FOREIGN KEY (setup_data_id) REFERENCES setup_data(id)
                                 );


                                 `,
	`                               CREATE TABLE IF NOT EXISTS telemetry_data (
                                   id INT NOT NULL AUTO_INCREMENT,
                                   session_uid BIGINT UNSIGNED NOT NULL,
                                   frame_identifier int NOT NULL,
																	 session_time DECIMAL(16,10),
                                   m_buttonStatus BIT(4),
                                   PRIMARY KEY (id),
                                   FOREIGN KEY (session_uid) REFERENCES race_event_directory(session_uid)
                                 );`,
	`                               CREATE TABLE IF NOT EXISTS car_telemetry_data (
                                   id INT NOT NULL AUTO_INCREMENT,
                                   telemetry_data_id INT NOT NULL,
																	 car_index INT NOT NULL,
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
                                   FOREIGN KEY (telemetry_data_id) REFERENCES telemetry_data(id)
                                 );


                                 `,
	`                               CREATE TABLE IF NOT EXISTS status_data (
                                   id INT NOT NULL AUTO_INCREMENT,
                                   session_uid BIGINT UNSIGNED NOT NULL,
                                   frame_identifier int NOT NULL,
																	 session_time DECIMAL(16,10),
                                   PRIMARY KEY (id),
                                   FOREIGN KEY (session_uid) REFERENCES race_event_directory(session_uid)
                                 );`,
	`                               CREATE TABLE IF NOT EXISTS car_status_data (
                                   id INT NOT NULL AUTO_INCREMENT,
                                   status_data_id INT NOT NULL,
																	 car_index INT NOT NULL,
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
                                   FOREIGN KEY (status_data_id) REFERENCES status_data(id)
                                 );`,
}

func deleteDatabase(db *sql.DB) error {
	_, err := db.Exec("DROP DATABASE F1_GO_MYSQL")
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

func createDatabaseTables(db *sql.DB) error {
	for i, stmt := range createTables {

		fmt.Print("   Create table ", tableNames[i], strings.Repeat(" ", (20-utf8.RuneCountInString(tableNames[i])))+"    ")
		if _, err := db.Exec("DESCRIBE " + tableNames[i]); err != nil {
			// MySQL error 1146 is "table does not exist"
			if mErr, ok := err.(*mysql.MySQLError); ok && mErr.Number == 1146 {

				_, err := db.Exec(stmt)
				if err != nil {
					color.Red("Error")
					db.Close()
					return err
				} else {
					color.Green("Success")
				}
			}
		} else {
			color.Yellow("Exists")
		}
	}

	// We are now finished making our tables
	// Close the connection and return success!
	db.Close()
	return nil
}

func add_race_event_directory_to_mysql(db *sql.DB, prepared_statement *sql.Stmt, packet structs.RaceEventDirectory) error {
	_, err = prepared_statement.Exec(
		packet.M_header.M_sessionUID,
		packet.M_header.M_packetFormat,
		packet.M_header.M_packetVersion,
		packet.M_header.M_playerCarIndex,
		packet.Session_start_time,
		packet.Session_end_time)
	if err != nil {
		fmt.Println("error adding race_event_directory to mysql, error:", err)
		return err
	}
	return nil
}

func add_motion_packet_to_mysql(db *sql.DB, prepared_statement *sql.Stmt, car_prepared_statement *sql.Stmt, packet structs.PacketMotionData) error {
	// First add motion_packet and get its id back
	res, err := prepared_statement.Exec(
		packet.M_header.M_sessionUID,
		packet.M_header.M_frameIdentifier,
		packet.M_header.M_sessionTime,
		packet.M_suspensionPosition[0],
		packet.M_suspensionPosition[1],
		packet.M_suspensionPosition[2],
		packet.M_suspensionPosition[3],
		packet.M_suspensionVelocity[0],
		packet.M_suspensionVelocity[1],
		packet.M_suspensionVelocity[2],
		packet.M_suspensionVelocity[3],
		packet.M_suspensionAcceleration[0],
		packet.M_suspensionAcceleration[1],
		packet.M_suspensionAcceleration[2],
		packet.M_suspensionAcceleration[3],
		packet.M_wheelSpeed[0],
		packet.M_wheelSpeed[1],
		packet.M_wheelSpeed[2],
		packet.M_wheelSpeed[3],
		packet.M_wheelSlip[0],
		packet.M_wheelSlip[1],
		packet.M_wheelSlip[2],
		packet.M_wheelSlip[3],
		packet.M_localVelocityX,
		packet.M_localVelocityY,
		packet.M_localVelocityZ,
		packet.M_angularVelocityX,
		packet.M_angularVelocityY,
		packet.M_angularVelocityZ,
		packet.M_angularAccelerationX,
		packet.M_angularAccelerationY,
		packet.M_angularAccelerationZ,
		packet.M_frontWheelsAngle)
	if err != nil {
		fmt.Println("error adding motion_packet to mysql, error:", err)
		return err
	} else {
		// If successfull, Get the id of the motion_packet
		id, err := res.LastInsertId()
		if err != nil {
			fmt.Println("error getting LastInsertId for motion_packet, error:", err)
			return err
		}

		// Loop through all the cars and add them to the MYSQL database
		for car_index, car := range packet.M_carMotionData {
			_, err = car_prepared_statement.Exec(
				id,
				car_index,
				car.M_worldPositionX,
				car.M_worldPositionY,
				car.M_worldPositionZ,
				car.M_worldVelocityX,
				car.M_worldVelocityY,
				car.M_worldVelocityZ,
				car.M_worldForwardDirX,
				car.M_worldForwardDirY,
				car.M_worldForwardDirZ,
				car.M_worldRightDirX,
				car.M_worldRightDirY,
				car.M_worldRightDirZ,
				car.M_gForceLateral,
				car.M_gForceLongitudinal,
				car.M_gForceVertical,
				car.M_yaw,
				car.M_pitch,
				car.M_roll)
			if err != nil {
				fmt.Println("error adding car_motion_packet to mysql, error:", err)
				return err
			}
		}
	}
	return nil
}

func add_session_packet_to_mysql(db *sql.DB, prepared_statement *sql.Stmt, car_prepared_statement *sql.Stmt, packet structs.PacketSessionData) error {
	// First add session_packet and get its id back
	res, err := prepared_statement.Exec(
		packet.M_header.M_sessionUID,
		packet.M_header.M_frameIdentifier,
		packet.M_header.M_sessionTime,
		packet.M_weather,
		packet.M_trackTemperature,
		packet.M_airTemperature,
		packet.M_totalLaps,
		packet.M_trackLength,
		packet.M_sessionType,
		packet.M_trackId,
		packet.M_era,
		packet.M_sessionTimeLeft,
		packet.M_sessionDuration,
		packet.M_pitSpeedLimit,
		packet.M_gamePaused,
		packet.M_isSpectating,
		packet.M_spectatorCarIndex,
		packet.M_sliProNativeSupport,
		packet.M_numMarshalZones,
		packet.M_safetyCarStatus,
		packet.M_networkGame)
	if err != nil {
		fmt.Println("error adding session_packet to mysql, error:", err)
		return err
	} else {
		// If successfull, Get the id of the session_packet
		id, err := res.LastInsertId()
		if err != nil {
			fmt.Println("error getting LastInsertId for session_packet, error:", err)
			return err
		}

		// Loop through all the cars and add them to the MYSQL database
		for car_index := 0; car_index < int(packet.M_numMarshalZones); car_index++ {
			_, err = car_prepared_statement.Exec(id, car_index, packet.M_marshalZones[car_index].M_zoneStart, packet.M_marshalZones[car_index].M_zoneFlag)
			if err != nil {
				fmt.Println("error adding car_session_packet to mysql, error:", err)
				return err
			}
		}
	}
	return nil
}

func add_lap_packet_to_mysql(db *sql.DB, prepared_statement *sql.Stmt, car_prepared_statement *sql.Stmt, packet structs.PacketLapData) error {
	// First add lap_packet and get its id back
	res, err := prepared_statement.Exec(
		packet.M_header.M_sessionUID,
		packet.M_header.M_frameIdentifier,
		packet.M_header.M_sessionTime)
	if err != nil {
		fmt.Println("error adding lap_packet to mysql, error:", err)
		return err
	} else {
		// If successfull, Get the id of the lap_packet
		id, err := res.LastInsertId()
		if err != nil {
			fmt.Println("error getting LastInsertId for lap_packet, error:", err)
			return err
		}

		// Loop through all the cars and add them to the MYSQL database
		for car_index, car := range packet.M_lapData {
			// fmt.Println(car.M_totalDistance)
			_, err = car_prepared_statement.Exec(
				id,
				car_index,
				car.M_lastLapTime,
				car.M_currentLapTime,
				car.M_bestLapTime,
				car.M_sector1Time,
				car.M_sector2Time,
				car.M_lapDistance,
				car.M_totalDistance,
				car.M_safetyCarDelta,
				car.M_carPosition,
				car.M_currentLapNum,
				car.M_pitStatus,
				car.M_sector,
				car.M_currentLapInvalid,
				car.M_penalties,
				car.M_gridPosition,
				car.M_driverStatus,
				car.M_resultStatus)
			if err != nil {
				fmt.Println("error adding car_lap_packet to mysql, error:", err)
				return err
			}
		}
	}
	return nil
}

func add_event_packet_to_mysql(db *sql.DB, prepared_statement *sql.Stmt, packet structs.PacketEventData) error {
	// First add lap_packet and get its id back
	_, err := prepared_statement.Exec(
		packet.M_header.M_sessionUID,
		packet.M_header.M_frameIdentifier,
		packet.M_header.M_sessionTime,
		packet.M_eventStringCode)
	if err != nil {
		fmt.Println("error adding event_packet to mysql, error:", err)
		return err
	}
	return nil
}

func add_participant_packet_to_mysql(db *sql.DB, prepared_statement *sql.Stmt, car_prepared_statement *sql.Stmt, packet structs.PacketParticipantsData) error {
	// First add lap_packet and get its id back
	res, err := prepared_statement.Exec(
		packet.M_header.M_sessionUID,
		packet.M_header.M_frameIdentifier,
		packet.M_header.M_sessionTime,
		packet.M_numCars)
	if err != nil {
		fmt.Println("error adding participant_packet to mysql, error:", err)
		return err
	} else {
		// If successfull, Get the id of the lap_packet
		id, err := res.LastInsertId()
		if err != nil {
			fmt.Println("error getting LastInsertId for participant_packet, error:", err)
			return err
		}

		// Loop through all the cars and add them to the MYSQL database
		for car_index, car := range packet.M_participants {
			_, err = car_prepared_statement.Exec(
				id,
				car_index,
				car.M_aiControlled,
				car.M_driverId,
				car.M_teamId,
				car.M_raceNumber,
				car.M_nationality,
				string(car.M_name[:]))
			if err != nil {
				fmt.Println("error adding car_participant_packet to mysql, error:", err)
				return err
			}
		}

	}
	return nil
}

func add_car_setup_packet_to_mysql(db *sql.DB, prepared_statement *sql.Stmt, car_prepared_statement *sql.Stmt, packet structs.PacketCarSetupData) error {
	// First add lap_packet and get its id back
	res, err := prepared_statement.Exec(
		packet.M_header.M_sessionUID,
		packet.M_header.M_frameIdentifier,
		packet.M_header.M_sessionTime)
	if err != nil {
		fmt.Println("error adding setup_packet to mysql, error:", err)
		return err
	} else {
		// If successfull, Get the id of the lap_packet
		id, err := res.LastInsertId()
		if err != nil {
			fmt.Println("error getting LastInsertId for setup_packet, error:", err)
			return err
		}

		// Loop through all the cars and add them to the MYSQL database
		for car_index, car := range packet.M_carSetups {
			_, err = car_prepared_statement.Exec(
				id,
				car_index,
				car.M_frontWing,
				car.M_rearWing,
				car.M_onThrottle,
				car.M_offThrottle,
				car.M_frontCamber,
				car.M_rearCamber,
				car.M_frontToe,
				car.M_rearToe,
				car.M_frontSuspension,
				car.M_rearSuspension,
				car.M_frontAntiRollBar,
				car.M_rearAntiRollBar,
				car.M_frontSuspensionHeight,
				car.M_rearSuspensionHeight,
				car.M_brakePressure,
				car.M_brakeBias,
				car.M_frontTyrePressure,
				car.M_rearTyrePressure,
				car.M_ballast,
				car.M_fuelLoad)
			if err != nil {
				fmt.Println("error adding car_setup_packet to mysql, error:", err)
				return err
			}
		}
	}
	return nil
}

func add_telemetry_packet_to_mysql(db *sql.DB, prepared_statement *sql.Stmt, car_prepared_statement *sql.Stmt, packet structs.PacketCarTelemetryData) error {
	// First add lap_packet and get its id back
	res, err := prepared_statement.Exec(
		packet.M_header.M_sessionUID,
		packet.M_header.M_frameIdentifier,
		packet.M_header.M_sessionTime,
		packet.M_buttonStatus)
	if err != nil {
		fmt.Println("error adding telemetry_packet to mysql, error:", err)
		return err
	} else {
		// If successfull, Get the id of the lap_packet
		id, err := res.LastInsertId()
		if err != nil {
			fmt.Println("error getting LastInsertId for telemetry_packet, error:", err)
			return err
		}

		// Loop through all the cars and add them to the MYSQL database
		for car_index, car := range packet.M_carTelemetryData {
			_, err = car_prepared_statement.Exec(
				id,
				car_index,
				car.M_speed,
				car.M_throttle,
				car.M_steer,
				car.M_brake,
				car.M_clutch,
				car.M_gear,
				car.M_engineRPM,
				car.M_drs,
				car.M_revLightsPercent,
				car.M_brakesTemperature[0],
				car.M_brakesTemperature[1],
				car.M_brakesTemperature[2],
				car.M_brakesTemperature[3],
				car.M_tyresSurfaceTemperature[0],
				car.M_tyresSurfaceTemperature[1],
				car.M_tyresSurfaceTemperature[2],
				car.M_tyresSurfaceTemperature[3],
				car.M_tyresInnerTemperature[0],
				car.M_tyresInnerTemperature[1],
				car.M_tyresInnerTemperature[2],
				car.M_tyresInnerTemperature[3],
				car.M_engineTemperature,
				car.M_tyresPressure[0],
				car.M_tyresPressure[1],
				car.M_tyresPressure[2],
				car.M_tyresPressure[3])
			if err != nil {
				fmt.Println("error car_telemetry_packet to mysql, error:", err)
				return err
			}
		}
	}
	return nil
}

func add_car_status_packet_to_mysql(db *sql.DB, prepared_statement *sql.Stmt, car_prepared_statement *sql.Stmt, packet structs.PacketCarStatusData) error {
	// First add lap_packet and get its id back
	res, err := prepared_statement.Exec(
		packet.M_header.M_sessionUID,
		packet.M_header.M_frameIdentifier,
		packet.M_header.M_sessionTime)
	if err != nil {
		fmt.Println("error adding status_packet to mysql, error:", err)
		return err
	} else {
		// If successfull, Get the id of the lap_packet
		id, err := res.LastInsertId()
		if err != nil {
			fmt.Println("error getting LastInsertId for status_packet, error:", err)
			return err
		}

		// Loop through all the cars and add them to the MYSQL database
		for car_index, car := range packet.M_carStatusData {
			_, err = car_prepared_statement.Exec(
				id,
				car_index,
				car.M_tractionControl,
				car.M_antiLockBrakes,
				car.M_fuelMix,
				car.M_frontBrakeBias,
				car.M_pitLimiterStatus,
				car.M_fuelInTank,
				car.M_fuelCapacity,
				car.M_maxRPM,
				car.M_idleRPM,
				car.M_maxGears,
				car.M_drsAllowed,
				car.M_tyresWear[0],
				car.M_tyresWear[1],
				car.M_tyresWear[2],
				car.M_tyresWear[3],
				car.M_tyreCompound,
				car.M_tyresDamage[0],
				car.M_tyresDamage[1],
				car.M_tyresDamage[2],
				car.M_tyresDamage[3],
				car.M_frontLeftWingDamage,
				car.M_frontRightWingDamage,
				car.M_rearWingDamage,
				car.M_engineDamage,
				car.M_gearBoxDamage,
				car.M_exhaustDamage,
				car.M_vehicleFiaFlags,
				car.M_ersStoreEnergy,
				car.M_ersDeployMode,
				car.M_ersHarvestedThisLapMGUK,
				car.M_ersHarvestedThisLapMGUH,
				car.M_ersDeployedThisLap)
			if err != nil {
				fmt.Println("error car_status_packet to mysql, error:", err)
				return err
			}
		}
	}

	return nil
}

func start_mysql() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Please enter your MYSQL password to connect to your MYSQL server:  ")
	fmt.Println("      user:      root")
	fmt.Print("      password:  ")
	scanner.Scan()
	mysql_password := scanner.Text()
	fmt.Print("\n")

	db, err := sql.Open("mysql", mysql_login_string_front+mysql_password+mysql_login_string_back)
	if err != nil {
		log.Println("mysql: could not get a connection: %v", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		log.Println("mysql: could not establish a good connection: %v", err)
		fmt.Println("Exiting...")
		os.Exit(1)
	} else {

		saved_mysql_password = mysql_login_string_front + mysql_password + mysql_login_string_back

		fmt.Print("Create F1_GO database  ")
		if _, err := db.Exec("USE F1_GO_MYSQL"); err != nil {
			// MySQL error 1049 is "database does not exist"
			if mErr, ok := err.(*mysql.MySQLError); ok && mErr.Number == 1049 {

				db, err = createDatabase(db)
				if err != nil {
					color.Red("Error")
					fmt.Println("Error creating F1_GO database!")
					log.Println(err)
					fmt.Println("Exiting...")
					os.Exit(1)
				} else {
					color.Green("Success")
				}
			}
		} else {
			color.Yellow("Exists")
		}

		err = createDatabaseTables(db)
		fmt.Print("F1_GO MYSQL tables     ")
		if err != nil {
			color.Red("Error")
			fmt.Println("Error creating F1_GO tables!")
			log.Println(err)
			fmt.Println("Exiting...")
			os.Exit(1)
		} else {
			color.Green("Done")
		}

		fmt.Print("F1_GO MYSQL database   ")
		color.Green("Done")
	}

	// Close the database connection
	db.Close()
}

func add_to_longterm_storage() {
	packets_to_add := true

	log.Println("mysql_login_string_front+saved_mysql_password+mysql_login_string_back:", saved_mysql_password)

	db, err := sql.Open("mysql", saved_mysql_password)
	if err != nil {
		log.Println("mysql: could not get a connection: %v", err)
	}

	if _, err := db.Exec("USE F1_GO_MYSQL"); err != nil {
		log.Println("mysql: error with statement 'USE F1_GO_MYSQL'", err)
	}

	// Defer the closing of the mysql database connection until we are finished with add_to_longterm_storage and return
	defer db.Close()

	if err := db.Ping(); err != nil {
		db.Close()
		log.Println("mysql: could not establish a good connection: %v", err)
		fmt.Println("Exiting...")
		os.Exit(1)
	} else {

		// Prepare statement for inserting data
		stmtIns_race_event_directory, err := db.Prepare("INSERT INTO race_event_directory (session_uid, M_packetFormat, packet_version, player_car_index, session_start, session_end) VALUES (?, ?, ?, ?, ?, ?)") // ? = placeholder
		if err != nil {
			// panic(err.Error()) // proper error handling instead of panic in your app
			log.Println("mysql: error with prepare statement stmtIns_race_event_directory")
		}

		// Prepare statement for inserting motion_data data
		stmtIns_motion_data, err := db.Prepare("INSERT INTO motion_data (session_uid, frame_identifier, session_time, suspension_position_rl, suspension_position_rr, suspension_position_fl, suspension_position_fr, M_suspensionVelocity_rl, M_suspensionVelocity_rr, M_suspensionVelocity_fl, M_suspensionVelocity_fr, suspension_acceleration_rl, suspension_acceleration_rr, suspension_acceleration_fl, suspension_acceleration_fr, wheel_speed_rl, wheel_speed_rr, wheel_speed_fl, wheel_speed_fr, wheel_slip_rl, wheel_slip_rr, wheel_slip_fl, wheel_slip_fr, local_velocity_x, local_velocity_y, local_velocity_z, angular_velocity_x, angular_velocity_y, angular_velocity_z, angular_acceleration_x, angular_acceleration_y, angular_acceleration_z, front_wheels_angle) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)") // ? = placeholder
		if err != nil {
			// panic(err.Error()) // proper error handling instead of panic in your app
			log.Println("mysql: error with prepare statement stmtIns_motion_data")
		}
		// Prepare statement for inserting car_motion_data data
		stmtIns_car_motion_data, err := db.Prepare("INSERT INTO car_motion_data (motion_packet_id, car_index, m_worldPositionX, m_worldPositionY, m_worldPositionZ, m_worldVelocityX, m_worldVelocityY, m_worldVelocityZ, m_worldForwardDirX, m_worldForwardDirY, m_worldForwardDirZ, m_worldRightDirX, m_worldRightDirY, m_worldRightDirZ, m_gForceLateral, m_gForceLongitudinal, m_gForceVertical, m_yaw, m_pitch, m_roll) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)") // ? = placeholder
		if err != nil {
			// panic(err.Error()) // proper error handling instead of panic in your app
			log.Println("mysql: error with prepare statement stmtIns_car_motion_data")
		}

		// Prepare statement for inserting session_data data
		stmtIns_session_data, err := db.Prepare("INSERT INTO session_data (session_uid, frame_identifier, session_time, m_weather, m_trackTemperature, m_airTemperature, m_totalLaps, m_trackLength, m_sessionType, m_trackId, m_era, m_sessionTimeLeft, m_sessionDuration, m_pitSpeedLimit, m_gamePaused, m_isSpectating, m_spectatorCarIndex, m_sliProNativeSupport, m_numMarshalZones, m_safetyCarStatus, m_networkGame) VALUES( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ? )") // ? = placeholder
		if err != nil {
			// panic(err.Error()) // proper error handling instead of panic in your app
			log.Println("mysql: error with prepare statement stmtIns_session_data")
		}
		// Prepare statement for inserting marshal_zone data
		stmtIns_marshal_zone, err := db.Prepare("INSERT INTO marshal_zone (session_data_id, car_index, m_zoneStart, m_zoneFlag) VALUES( ?, ?, ?, ? )") // ? = placeholder
		if err != nil {
			// panic(err.Error()) // proper error handling instead of panic in your app
			log.Println("mysql: error with prepare statement stmtIns_marshal_zone")
		}

		// Prepare statement for inserting lap_data data
		stmtIns_lap_data, err := db.Prepare("INSERT INTO lap_data (session_uid, frame_identifier, session_time) VALUES( ?, ?, ? )") // ? = placeholder
		if err != nil {
			// panic(err.Error()) // proper error handling instead of panic in your app
			log.Println("mysql: error with prepare statement stmtIns_lap_data")
		}
		// Prepare statement for inserting car_lap_data data
		stmtIns_car_lap_data, err := db.Prepare("INSERT INTO car_lap_data (lap_data_id, car_index, m_lastLapTime, m_currentLapTime, m_bestLapTime, m_sector1Time, m_sector2Time, m_lapDistance, m_totalDistance, m_safetyCarDelta, m_carPosition, m_currentLapNum, m_pitStatus, m_sector, m_currentLapInvalid, m_penalties, m_gridPosition, m_driverStatus, m_resultStatus) VALUES( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ? )") // ? = placeholder
		if err != nil {
			// panic(err.Error()) // proper error handling instead of panic in your app
			log.Println("mysql: error with prepare statement stmtIns_car_lap_data")
		}

		// Prepare statement for inserting event_data data
		stmtIns_event_data, err := db.Prepare("INSERT INTO event_data (session_uid, frame_identifier, session_time, m_eventStringCode) VALUES( ?, ?, ?, ? )") // ? = placeholder
		if err != nil {
			// panic(err.Error()) // proper error handling instead of panic in your app
			log.Println("mysql: error with prepare statement stmtIns_event_data")
		}

		// Prepare statement for inserting participant_data data
		stmtIns_participant_data, err := db.Prepare("INSERT INTO participant_data (session_uid, frame_identifier, session_time, m_numCars) VALUES( ?, ?, ?, ? )") // ? = placeholder
		if err != nil {
			// panic(err.Error()) // proper error handling instead of panic in your app
			log.Println("mysql: error with prepare statement stmtIns_participant_data")
		}
		// Prepare statement for inserting car_participant_data data
		stmtIns_car_participant_data, err := db.Prepare("INSERT INTO car_participant_data (participant_data_id, car_index, m_aiControlled, m_driverId, m_teamId, m_raceNumber, m_nationality, m_name) VALUES( ?, ?, ?, ?, ?, ?, ?, ? )") // ? = placeholder
		if err != nil {
			// panic(err.Error()) // proper error handling instead of panic in your app
			log.Println("mysql: error with prepare statement stmtIns_car_participant_data")
		}

		// Prepare statement for inserting setup_data data
		stmtIns_setup_data, err := db.Prepare("INSERT INTO setup_data (session_uid, frame_identifier, session_time) VALUES( ?, ?, ? )") // ? = placeholder
		if err != nil {
			// panic(err.Error()) // proper error handling instead of panic in your app
			log.Println("mysql: error with prepare statement stmtIns_setup_data")
		}
		// Prepare statement for inserting car_setup_data data
		stmtIns_car_setup_data, err := db.Prepare("INSERT INTO car_setup_data (setup_data_id, car_index, m_frontWing, m_rearWing, m_onThrottle, m_offThrottle, m_frontCamber, m_rearCamber, m_frontToe, m_rearToe, m_frontSuspension, m_rearSuspension, m_frontAntiRollBar, m_rearAntiRollBar, m_frontSuspensionHeight, m_rearSuspensionHeight, m_brakePressure, m_brakeBias, m_frontTyrePressure, m_rearTyrePressure, m_ballast, m_fuelLoad) VALUES( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ? )") // ? = placeholder
		if err != nil {
			// panic(err.Error()) // proper error handling instead of panic in your app
			log.Println("mysql: error with prepare statement stmtIns_car_setup_data")
		}

		// Prepare statement for inserting telemetry data
		stmtIns_telemetry_data, err := db.Prepare("INSERT INTO telemetry_data (session_uid, frame_identifier, session_time, m_buttonStatus) VALUES( ?, ?, ?, ? )") // ? = placeholder
		if err != nil {
			// panic(err.Error()) // proper error handling instead of panic in your app
			log.Println("mysql: error with prepare statement stmtIns_telemetry_data")
		}
		// Prepare statement for inserting car_setup_data data
		stmtIns_car_telemetry_data, err := db.Prepare("INSERT INTO car_telemetry_data (telemetry_data_id, car_index, m_speed, m_throttle, m_steer, m_brake, m_clutch, m_gear, m_engineRPM, m_drs, m_revLightsPercent, m_brakesTemperature_rl, m_brakesTemperature_rr, m_brakesTemperature_fl, m_brakesTemperature_fr, m_tyresSurfaceTemperature_rl, m_tyresSurfaceTemperature_rr, m_tyresSurfaceTemperature_fl, m_tyresSurfaceTemperature_fr, m_tyresInnerTemperature_rl, m_tyresInnerTemperature_rr, m_tyresInnerTemperature_fl, m_tyresInnerTemperature_fr, m_engineTemperature, m_tyresPressure_rl, m_tyresPressure_rr, m_tyresPressure_fl, m_tyresPressure_fr) VALUES( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ? )") // ? = placeholder
		if err != nil {
			// panic(err.Error()) // proper error handling instead of panic in your app
			log.Println("mysql: error with prepare statement stmtIns_car_telemetry_data")
		}

		// Prepare statement for inserting status_data data
		stmtIns_status_data, err := db.Prepare("INSERT INTO status_data (session_uid, frame_identifier, session_time) VALUES( ?, ?, ? )") // ? = placeholder
		if err != nil {
			// panic(err.Error()) // proper error handling instead of panic in your app
			log.Println("mysql: error with prepare statement stmtIns_status_data")
		}

		// Prepare statement for inserting car_status_data data
		stmtIns_car_status_data, err := db.Prepare("INSERT INTO car_status_data (status_data_id, car_index, m_tractionControl, m_antiLockBrakes, m_fuelMix, m_frontBrakeBias, m_pitLimiterStatus, m_fuelInTank, m_fuelCapacity, m_maxRPM, m_idleRPM, m_maxGears, m_drsAllowed, m_tyresWear_rl, m_tyresWear_rr, m_tyresWear_fl, m_tyresWear_fr, m_tyreCompound, m_tyresDamage_rl, m_tyresDamage_rr, m_tyresDamage_fl, m_tyresDamage_fr, m_frontLeftWingDamage, m_frontRightWingDamage, m_rearWingDamage, m_engineDamage, m_gearBoxDamage, m_exhaustDamage, m_vehicleFiaFlags, m_ersStoreEnergy, m_ersDeployMode, m_ersHarvestedThisLapMGUK, m_ersHarvestedThisLapMGUH, m_ersDeployedThisLap) VALUES( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,?, ?, ?, ? )") // ? = placeholder
		if err != nil {
			// panic(err.Error()) // proper error handling instead of panic in your app
			log.Println("mysql: error with prepare statement stmtIns_car_status_data")
		}

		defer stmtIns_race_event_directory.Close() // Close the statement when we leave main() / the program terminates

		defer stmtIns_motion_data.Close()     // Close the statement when we leave main() / the program terminates
		defer stmtIns_car_motion_data.Close() // Close the statement when we leave main() / the program terminates

		defer stmtIns_session_data.Close() // Close the statement when we leave main() / the program terminates
		defer stmtIns_marshal_zone.Close() // Close the statement when we leave main() / the program terminates

		defer stmtIns_lap_data.Close()     // Close the statement when we leave main() / the program terminates
		defer stmtIns_car_lap_data.Close() // Close the statement when we leave main() / the program terminates

		defer stmtIns_event_data.Close() // Close the statement when we leave main() / the program terminates

		defer stmtIns_participant_data.Close()     // Close the statement when we leave main() / the program terminates
		defer stmtIns_car_participant_data.Close() // Close the statement when we leave main() / the program terminates

		defer stmtIns_setup_data.Close()     // Close the statement when we leave main() / the program terminates
		defer stmtIns_car_setup_data.Close() // Close the statement when we leave main() / the program terminates

		defer stmtIns_telemetry_data.Close()     // Close the statement when we leave main() / the program terminates
		defer stmtIns_car_telemetry_data.Close() // Close the statement when we leave main() / the program terminates

		defer stmtIns_status_data.Close()     // Close the statement when we leave main() / the program terminates
		defer stmtIns_car_status_data.Close() // Close the statement when we leave main() / the program terminates

		for packets_to_add {
			select {
			case motion_packet := <-atm_motion_packet:
				// fmt.Println(motion_packet, "atm_motion_packet")
				if err := add_motion_packet_to_mysql(db, stmtIns_motion_data, stmtIns_car_motion_data, motion_packet); err != nil {
					log.Println("add_to_longterm_storage: error adding motion_packet to mysql: %v", err)
					log.Println("motion_packet.M_header.M_sessionUID:", motion_packet.M_header.M_sessionUID)
				}

			case session_packet := <-atm_session_packet:
				// fmt.Println(session_packet, "atm_session_packet")
				if err := add_session_packet_to_mysql(db, stmtIns_session_data, stmtIns_marshal_zone, session_packet); err != nil {
					log.Println("add_to_longterm_storage: error adding session_packet to mysql: %v", err)
					log.Println("session_packet.M_header.M_sessionUID:", session_packet.M_header.M_sessionUID)
				}

			case lap_packet := <-atm_lap_packet:
				// fmt.Println(motion_packet, "atm_lap_packet")
				if err := add_lap_packet_to_mysql(db, stmtIns_lap_data, stmtIns_car_lap_data, lap_packet); err != nil {
					log.Println("add_to_longterm_storage: error adding lap_packet to mysql: %v", err)
					log.Println("lap_packet.M_header.M_sessionUID:", lap_packet.M_header.M_sessionUID)
				}

			case event_packet := <-atm_event_packet:
				// fmt.Println(event_packet, "atm_event_packet")
				if err := add_event_packet_to_mysql(db, stmtIns_event_data, event_packet); err != nil {
					log.Println("add_to_longterm_storage: error adding event_packet to mysql: %v", err)
					log.Println("event_packet.M_header.M_sessionUID:", event_packet.M_header.M_sessionUID)
				}

			case participant_packet := <-atm_participant_packet:
				// fmt.Println(participant_packet, "atm_participant_packet")
				if err := add_participant_packet_to_mysql(db, stmtIns_participant_data, stmtIns_car_participant_data, participant_packet); err != nil {
					log.Println("add_to_longterm_storage: error adding participant_packet to mysql: %v", err)
					log.Println("participant_packet.M_header.M_sessionUID:", participant_packet.M_header.M_sessionUID)
				}

			case car_setup_packet := <-atm_car_setup_packet:
				// fmt.Println(car_setup_packet, "atm_car_setup_packet")
				if err := add_car_setup_packet_to_mysql(db, stmtIns_setup_data, stmtIns_car_setup_data, car_setup_packet); err != nil {
					log.Println("add_to_longterm_storage: error adding car_setup_packet to mysql: %v", err)
					log.Println("car_setup_packet.M_header.M_sessionUID:", car_setup_packet.M_header.M_sessionUID)
				}

			case telemetry_packet := <-atm_telemetry_packet:
				// fmt.Println(telemetry_packet, "atm_telemetry_packet")
				if err := add_telemetry_packet_to_mysql(db, stmtIns_telemetry_data, stmtIns_car_telemetry_data, telemetry_packet); err != nil {
					log.Println("add_to_longterm_storage: error adding telemetry_packet to mysql: %v", err)
					log.Println("telemetry_packet.M_header.M_sessionUID:", telemetry_packet.M_header.M_sessionUID)
				}

			case car_status_packet := <-atm_car_status_packet:
				// fmt.Println(car_status_packet, "atm_car_status_packet")
				if err := add_car_status_packet_to_mysql(db, stmtIns_status_data, stmtIns_car_status_data, car_status_packet); err != nil {
					log.Println("add_to_longterm_storage: error adding car_status_packet to mysql: %v", err)
					log.Println("car_status_packet.M_header.M_sessionUID:", car_status_packet.M_header.M_sessionUID)
				}

			case race_event_directory_data := <-atm_race_event_directory:
				// fmt.Println(car_status_packet, "atm_car_status_packet")
				if err := add_race_event_directory_to_mysql(db, stmtIns_race_event_directory, race_event_directory_data); err != nil {
					log.Println("add_to_longterm_storage: error adding race_event_directory_data to mysql: %v", err)
					log.Println("race_event_directory_data.M_header.M_sessionUID:", race_event_directory_data.M_header.M_sessionUID)
				}

			case _ = <-redis_done:
				fmt.Println("Redis finished sending data to MYSQL")
				packets_to_add = false

			}
		}
	}

	return

}
