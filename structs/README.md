## PACKET TYPES

Taken from codemasters forum <https://forums.codemasters.com/discussion/136948/f1-2018-udp-specification>

The main change for 2018 is the introduction of multiple packet types: each packet can now carry different types of data rather than having one packet which contains everything. A header has been added to each packet as well so that versioning can be tracked and it will be easier for applications to check they are interpreting the incoming data in the correct way.

Each packet has the following header:

```
struct PacketHeader
{
    uint16    m_packetFormat;         // 2018
    uint8     m_packetVersion;        // Version of this packet type, all start from 1
    uint8     m_packetId;             // Identifier for the packet type, see below
    uint64    m_sessionUID;           // Unique identifier for the session
    float     m_sessionTime;          // Session timestamp
    uint      m_frameIdentifier;      // Identifier for the frame the data was retrieved on
    uint8     m_playerCarIndex;       // Index of player's car in the array
};
```

## MOTION PACKET

> The motion packet gives physics data for all the cars being driven. There is additional data for the car being driven with the goal of being able to drive a motion platform setup.

N.B. For the normalised vectors below, to convert to float values divide by 32767.0f. 16-bit signed values are used to pack the data and on the assumption that direction values are always between -1.0f and 1.0f.

Frequency: Rate as specified in menus

Size: 1341 bytes

```
struct CarMotionData
{
    float         m_worldPositionX;           // World space X position
    float         m_worldPositionY;           // World space Y position
    float         m_worldPositionZ;           // World space Z position
    float         m_worldVelocityX;           // Velocity in world space X
    float         m_worldVelocityY;           // Velocity in world space Y
    float         m_worldVelocityZ;           // Velocity in world space Z
    int16         m_worldForwardDirX;         // World space forward X direction (normalised)
    int16         m_worldForwardDirY;         // World space forward Y direction (normalised)
    int16         m_worldForwardDirZ;         // World space forward Z direction (normalised)
    int16         m_worldRightDirX;           // World space right X direction (normalised)
    int16         m_worldRightDirY;           // World space right Y direction (normalised)
    int16         m_worldRightDirZ;           // World space right Z direction (normalised)
    float         m_gForceLateral;            // Lateral G-Force component
    float         m_gForceLongitudinal;       // Longitudinal G-Force component
    float         m_gForceVertical;           // Vertical G-Force component
    float         m_yaw;                      // Yaw angle in radians
    float         m_pitch;                    // Pitch angle in radians
    float         m_roll;                     // Roll angle in radians
};
```

```
struct PacketMotionData
{
    PacketHeader    m_header;               // Header

    CarMotionData   m_carMotionData[20];    // Data for all cars on track

    // Extra player car ONLY data
    float         m_suspensionPosition[4];       // Note: All wheel arrays have the following order:
    float         m_suspensionVelocity[4];       // RL, RR, FL, FR
    float         m_suspensionAcceleration[4];   // RL, RR, FL, FR
    float         m_wheelSpeed[4];               // Speed of each wheel
    float         m_wheelSlip[4];                // Slip ratio for each wheel
    float         m_localVelocityX;              // Velocity in local space
    float         m_localVelocityY;              // Velocity in local space
    float         m_localVelocityZ;              // Velocity in local space
    float         m_angularVelocityX;            // Angular velocity x-component
    float         m_angularVelocityY;            // Angular velocity y-component
    float         m_angularVelocityZ;            // Angular velocity z-component
    float         m_angularAccelerationX;        // Angular velocity x-component
    float         m_angularAccelerationY;        // Angular velocity y-component
    float         m_angularAccelerationZ;        // Angular velocity z-component
    float         m_frontWheelsAngle;            // Current front wheels angle in radians
};
```

## SESSION PACKET

> The session packet includes details about the current session in progress.

Frequency: 2 per second

Size: 147 bytes

```
struct MarshalZone
{
    float  m_zoneStart;   // Fraction (0..1) of way through the lap the marshal zone starts
    int8   m_zoneFlag;    // -1 = invalid/unknown, 0 = none, 1 = green, 2 = blue, 3 = yellow, 4 = red
};
```

```
struct PacketSessionData
{
    PacketHeader    m_header;                   // Header

    uint8           m_weather;                  // Weather - 0 = clear, 1 = light cloud, 2 = overcast
                                                // 3 = light rain, 4 = heavy rain, 5 = storm
    int8        m_trackTemperature;        // Track temp. in degrees celsius
    int8        m_airTemperature;          // Air temp. in degrees celsius
    uint8           m_totalLaps;               // Total number of laps in this race
    uint16          m_trackLength;               // Track length in metres
    uint8           m_sessionType;             // 0 = unknown, 1 = P1, 2 = P2, 3 = P3, 4 = Short P
                                                // 5 = Q1, 6 = Q2, 7 = Q3, 8 = Short Q, 9 = OSQ
                                                // 10 = R, 11 = R2, 12 = Time Trial
    int8            m_trackId;                 // -1 for unknown, 0-21 for tracks, see appendix
    uint8           m_era;                      // Era, 0 = modern, 1 = classic
    uint16          m_sessionTimeLeft;        // Time left in session in seconds
    uint16          m_sessionDuration;         // Session duration in seconds
    uint8           m_pitSpeedLimit;          // Pit speed limit in kilometres per hour
    uint8           m_gamePaused;               // Whether the game is paused
    uint8           m_isSpectating;            // Whether the player is spectating
    uint8           m_spectatorCarIndex;      // Index of the car being spectated
    uint8           m_sliProNativeSupport;    // SLI Pro support, 0 = inactive, 1 = active
    uint8           m_numMarshalZones;             // Number of marshal zones to follow
    MarshalZone     m_marshalZones[21];         // List of marshal zones – max 21
    uint8           m_safetyCarStatus;          // 0 = no safety car, 1 = full safety car
                                                // 2 = virtual safety car
    uint8          m_networkGame;              // 0 = offline, 1 = online
};
```

## LAP DATA PACKET

> The lap data packet gives details of all the cars in the session.

Frequency: Rate as specified in menus

Size: 841 bytes

```
struct LapData
{
    float       m_lastLapTime;           // Last lap time in seconds
    float       m_currentLapTime;        // Current time around the lap in seconds
    float       m_bestLapTime;           // Best lap time of the session in seconds
    float       m_sector1Time;           // Sector 1 time in seconds
    float       m_sector2Time;           // Sector 2 time in seconds
    float       m_lapDistance;           // Distance vehicle is around current lap in metres – could
                                         // be negative if line hasn’t been crossed yet
    float       m_totalDistance;         // Total distance travelled in session in metres – could
                                         // be negative if line hasn’t been crossed yet
    float       m_safetyCarDelta;        // Delta in seconds for safety car
    uint8       m_carPosition;           // Car race position
    uint8       m_currentLapNum;         // Current lap number
    uint8       m_pitStatus;             // 0 = none, 1 = pitting, 2 = in pit area
    uint8       m_sector;                // 0 = sector1, 1 = sector2, 2 = sector3
    uint8       m_currentLapInvalid;     // Current lap invalid - 0 = valid, 1 = invalid
    uint8       m_penalties;             // Accumulated time penalties in seconds to be added
    uint8       m_gridPosition;          // Grid position the vehicle started the race in
    uint8       m_driverStatus;          // Status of driver - 0 = in garage, 1 = flying lap
                                         // 2 = in lap, 3 = out lap, 4 = on track
    uint8       m_resultStatus;          // Result status - 0 = invalid, 1 = inactive, 2 = active
                                         // 3 = finished, 4 = disqualified, 5 = not classified
                                         // 6 = retired
};
```

```
struct PacketLapData
{
    PacketHeader    m_header;              // Header

    LapData         m_lapData[20];         // Lap data for all cars on track
};
```

## EVENT PACKET

> This packet gives details of events that happen during the course of the race.

Frequency: When the event occurs

Size: 25 bytes

```
struct PacketEventData
{
    PacketHeader    m_header;               // Header

    uint8           m_eventStringCode[4];   // Event string code, see above
};
```

## PARTICIPANTS PACKET

> This is a list of participants in the race. If the vehicle is controlled by AI, then the name will be the driver name. If this is a multiplayer game, the names will be the Steam Id on PC, or the LAN name if appropriate. On Xbox One, the names will always be the driver name, on PS4 the name will be the LAN name if playing a LAN game, otherwise it will be the driver name.

Frequency: Every 5 seconds

Size: 1082 bytes

```
struct ParticipantData
{
    uint8      m_aiControlled;           // Whether the vehicle is AI (1) or Human (0) controlled
    uint8      m_driverId;               // Driver id - see appendix
    uint8      m_teamId;                 // Team id - see appendix
    uint8      m_raceNumber;             // Race number of the car
    uint8      m_nationality;            // Nationality of the driver
    char       m_name[48];               // Name of participant in UTF-8 format – null terminated
                                         // Will be truncated with … (U+2026) if too long
};
```

```
struct PacketParticipantsData
{
    PacketHeader    m_header;            // Header

    uint8           m_numCars;           // Number of cars in the data
    ParticipantData m_participants[20];
};
```

## CAR SETUPS PACKET

> This packet details the car setups for each vehicle in the session. Note that in multiplayer games, other player cars will appear as blank, you will only be able to see your car setup and AI cars.

Frequency: Every 5 seconds

Size: 841 bytes

```
struct CarSetupData
{
    uint8     m_frontWing;                // Front wing aero
    uint8     m_rearWing;                 // Rear wing aero
    uint8     m_onThrottle;               // Differential adjustment on throttle (percentage)
    uint8     m_offThrottle;              // Differential adjustment off throttle (percentage)
    float     m_frontCamber;              // Front camber angle (suspension geometry)
    float     m_rearCamber;               // Rear camber angle (suspension geometry)
    float     m_frontToe;                 // Front toe angle (suspension geometry)
    float     m_rearToe;                  // Rear toe angle (suspension geometry)
    uint8     m_frontSuspension;          // Front suspension
    uint8     m_rearSuspension;           // Rear suspension
    uint8     m_frontAntiRollBar;         // Front anti-roll bar
    uint8     m_rearAntiRollBar;          // Front anti-roll bar
    uint8     m_frontSuspensionHeight;    // Front ride height
    uint8     m_rearSuspensionHeight;     // Rear ride height
    uint8     m_brakePressure;            // Brake pressure (percentage)
    uint8     m_brakeBias;                // Brake bias (percentage)
    float     m_frontTyrePressure;        // Front tyre pressure (PSI)
    float     m_rearTyrePressure;         // Rear tyre pressure (PSI)
    uint8     m_ballast;                  // Ballast
    float     m_fuelLoad;                 // Fuel load
};
```

```
struct PacketCarSetupData
{
    PacketHeader    m_header;            // Header

    CarSetupData    m_carSetups[20];
};
```

## CAR TELEMETRY PACKET

> This packet details telemetry for all the cars in the race. It details various values that would be recorded on the car such as speed, throttle application, DRS etc.

Frequency: Rate as specified in menus

Size: 1085 bytes

```
struct CarTelemetryData
{
    uint16    m_speed;                      // Speed of car in kilometres per hour
    uint8     m_throttle;                   // Amount of throttle applied (0 to 100)
    int8      m_steer;                      // Steering (-100 (full lock left) to 100 (full lock right))
    uint8     m_brake;                      // Amount of brake applied (0 to 100)
    uint8     m_clutch;                     // Amount of clutch applied (0 to 100)
    int8      m_gear;                       // Gear selected (1-8, N=0, R=-1)
    uint16    m_engineRPM;                  // Engine RPM
    uint8     m_drs;                        // 0 = off, 1 = on
    uint8     m_revLightsPercent;           // Rev lights indicator (percentage)
    uint16    m_brakesTemperature[4];       // Brakes temperature (celsius)
    uint16    m_tyresSurfaceTemperature[4]; // Tyres surface temperature (celsius)
    uint16    m_tyresInnerTemperature[4];   // Tyres inner temperature (celsius)
    uint16    m_engineTemperature;          // Engine temperature (celsius)
    float     m_tyresPressure[4];           // Tyres pressure (PSI)
};
```

```
struct PacketCarTelemetryData
{
    PacketHeader        m_header;                // Header

    CarTelemetryData    m_carTelemetryData[20];

    uint32              m_buttonStatus;         // Bit flags specifying which buttons are being
                                                // pressed currently - see appendices
};
```

## CAR STATUS PACKET

> This packet details car statuses for all the cars in the race. It includes values such as the damage readings on the car.

Frequency: 2 per second

Size: 1061 bytes

```
struct CarStatusData
{
    uint8       m_tractionControl;          // 0 (off) - 2 (high)
    uint8       m_antiLockBrakes;           // 0 (off) - 1 (on)
    uint8       m_fuelMix;                  // Fuel mix - 0 = lean, 1 = standard, 2 = rich, 3 = max
    uint8       m_frontBrakeBias;           // Front brake bias (percentage)
    uint8       m_pitLimiterStatus;         // Pit limiter status - 0 = off, 1 = on
    float       m_fuelInTank;               // Current fuel mass
    float       m_fuelCapacity;             // Fuel capacity
    uint16      m_maxRPM;                   // Cars max RPM, point of rev limiter
    uint16      m_idleRPM;                  // Cars idle RPM
    uint8       m_maxGears;                 // Maximum number of gears
    uint8       m_drsAllowed;               // 0 = not allowed, 1 = allowed, -1 = unknown
    uint8       m_tyresWear[4];             // Tyre wear percentage
    uint8       m_tyreCompound;             // Modern - 0 = hyper soft, 1 = ultra soft
                                            // 2 = super soft, 3 = soft, 4 = medium, 5 = hard
                                            // 6 = super hard, 7 = inter, 8 = wet
                                            // Classic - 0-6 = dry, 7-8 = wet
    uint8       m_tyresDamage[4];           // Tyre damage (percentage)
    uint8       m_frontLeftWingDamage;      // Front left wing damage (percentage)
    uint8       m_frontRightWingDamage;     // Front right wing damage (percentage)
    uint8       m_rearWingDamage;           // Rear wing damage (percentage)
    uint8       m_engineDamage;             // Engine damage (percentage)
    uint8       m_gearBoxDamage;            // Gear box damage (percentage)
    uint8       m_exhaustDamage;            // Exhaust damage (percentage)
    int8        m_vehicleFiaFlags;          // -1 = invalid/unknown, 0 = none, 1 = green
                                            // 2 = blue, 3 = yellow, 4 = red
    float       m_ersStoreEnergy;           // ERS energy store in Joules
    uint8       m_ersDeployMode;            // ERS deployment mode, 0 = none, 1 = low, 2 = medium
                                            // 3 = high, 4 = overtake, 5 = hotlap
    float       m_ersHarvestedThisLapMGUK;  // ERS energy harvested this lap by MGU-K
    float       m_ersHarvestedThisLapMGUH;  // ERS energy harvested this lap by MGU-H
    float       m_ersDeployedThisLap;       // ERS energy deployed this lap
};
```

```
struct PacketCarStatusData
{
    PacketHeader        m_header;            // Header

    CarStatusData       m_carStatusData[20];
};
```

# Packet elements sizes in bytes

> All sizes are in bytes

## HEADER

> Size: 21 bytes

```
          type PacketHeader struct {
2             M_packetFormat    uint16  // 2018
1             M_packetVersion   uint8   // Version of this packet type, all start from 1
1             M_packetId        uint8   // Identifier for the packet type, see below
8             M_sessionUID      uint64  // Unique identifier for the session
4             M_sessionTime     float32 // Session timestamp
4             M_frameIdentifier uint32  // Identifier for the frame the data was retrieved on
1             M_playerCarIndex  uint8   // Index of player's car in the array
          }

Total: 20
```

## MOTION PACKET

> Size: 1341 bytes

```
          type CarMotionData struct {
4             M_worldPositionX     float32 // World space X position
4             M_worldPositionY     float32 // World space Y position
4             M_worldPositionZ     float32 // World space Z position
4             M_worldVelocityX     float32 // Velocity in world space X
4             M_worldVelocityY     float32 // Velocity in world space Y
4             M_worldVelocityZ     float32 // Velocity in world space Z
2             M_worldForwardDirX   int16   // World space forward X direction (normalised)
2             M_worldForwardDirY   int16   // World space forward Y direction (normalised)
2             M_worldForwardDirZ   int16   // World space forward Z direction (normalised)
2             M_worldRightDirX     int16   // World space right X direction (normalised)
2             M_worldRightDirY     int16   // World space right Y direction (normalised)
2             M_worldRightDirZ     int16   // World space right Z direction (normalised)
4             M_gForceLateral      float32 // Lateral G-Force component
4             M_gForceLongitudinal float32 // Longitudinal G-Force component
4             M_gForceVertical     float32 // Vertical G-Force component
4             M_yaw                float32 // Yaw angle in radians
4             M_pitch              float32 // Pitch angle in radians
4             M_roll               float32 // Roll angle in radians
          }
Total: 60

          type PacketMotionData struct {
21            M_header PacketHeader // Header

1,200       M_carMotionData [20]CarMotionData // Data for all cars on track

              // Extra player car ONLY data
16             M_suspensionPosition     [4]float32 // Note: All wheel arrays have the following order:
16             M_suspensionVelocity     [4]float32 // RL, RR, FL, FR
16            M_suspensionAcceleration [4]float32 // RL, RR, FL, FR
16            M_wheelSpeed             [4]float32 // Speed of each wheel
16            M_wheelSlip              [4]float32 // Slip ratio for each wheel
4             M_localVelocityX         float32    // Velocity in local space
4             M_localVelocityY         float32    // Velocity in local space
4             M_localVelocityZ         float32    // Velocity in local space
4             M_angularVelocityX       float32    // Angular velocity x-component
4             M_angularVelocityY       float32    // Angular velocity y-component
4             M_angularVelocityZ       float32    // Angular velocity z-component
4             M_angularAccelerationX   float32    // Angular velocity x-component
4             M_angularAccelerationY   float32    // Angular velocity y-component
4             M_angularAccelerationZ   float32    // Angular velocity z-component
4             M_frontWheelsAngle       float32    // Current front wheels angle in radians
        }
Total: 1341
```

## SESION PACKET

> Size: 147 bytes

```
          type MarshalZone struct {
4             M_zoneStart float32 // Fraction (0..1) of way through the lap the marshal zone starts
1             M_zoneFlag  int8    // -1 = invalid/unknown, 0 = none, 1 = green, 2 = blue, 3 = yellow, 4 = red
          }
Total: 5

          type PacketSessionData struct {
21            M_header PacketHeader // Header

1             M_weather             uint8           // Weather - 0 = clear, 1 = light cloud, 2 = overcast, 3 = light rain, 4 = heavy rain, 5 = storm
1             M_trackTemperature    int8            // Track temp. in degrees celsius
1             M_airTemperature      int8            // Air temp. in degrees celsius
1             M_totalLaps           uint8           // Total number of laps in this race
2             M_trackLength         uint16          // Track length in metre
1             M_sessionType         uint8           // 0 = unknown, 1 = P1, 2 = P2, 3 = P3, 4 = Short P, 5 = Q1, 6 = Q2, 7 = Q3, 8 = Short Q, 9 = OSQ, 10 = R, 11 = R2, 12 = Time Trial
1             M_trackId             int8            // -1 for unknown, 0-21 for tracks, see appendix
1             M_era                 uint8           // Era, 0 = modern, 1 = classic
2             M_sessionTimeLeft     uint16          // Time left in session in seconds
2             M_sessionDuration     uint16          // Session duration in seconds
1             M_pitSpeedLimit       uint8           // Pit speed limit in kilometres per hour
1             M_gamePaused          uint8           // Whether the game is paused
1             M_isSpectating        uint8           // Whether the player is spectating
1             M_spectatorCarIndex   uint8           // Index of the car being spectated
1             M_sliProNativeSupport uint8           // SLI Pro support, 0 = inactive, 1 = active
1             M_numMarshalZones     uint8           // Number of marshal zones to follow
105            M_marshalZones        [21]MarshalZone // List of marshal zones – max 21
1             M_safetyCarStatus     uint8           // 0 = no safety car, 1 = full safety car, 2 = virtual safety car
1             M_networkGame         uint8           // 0 = offline, 1 = online
        }
Total: 147
```

## LAP DATA PACKET

> Size: 841 bytes

```
          type LapData struct {
4             M_lastLapTime       float32 // Last lap time in seconds
4             M_currentLapTime    float32 // Current time around the lap in seconds
4             M_bestLapTime       float32 // Best lap time of the session in seconds
4             M_sector1Time       float32 // Sector 1 time in seconds
4             M_sector2Time       float32 // Sector 2 time in seconds
4             M_lapDistance       float32 // Distance vehicle is around current lap in metres – could be negative if line hasn’t been crossed yet
4             M_totalDistance     float32 // Total distance travelled in session in metres – could be negative if line hasn’t been crossed yet
4             M_safetyCarDelta    float32 // Delta in seconds for safety car
1             M_carPosition       uint8   // Car race position
1             M_currentLapNum     uint8   // Current lap number
1             M_pitStatus         uint8   // 0 = none, 1 = pitting, 2 = in pit area
1             M_sector            uint8   // 0 = sector1, 1 = sector2, 2 = sector3
1             M_currentLapInvalid uint8   // Current lap invalid - 0 = valid, 1 = invalid
1             M_penalties         uint8   // Accumulated time penalties in seconds to be added
1             M_gridPosition      uint8   // Grid position the vehicle started the race in
1             M_driverStatus      uint8   // Status of driver - 0 = in garage, 1 = flying lap, 2 = in lap, 3 = out lap, 4 = on track
1             M_resultStatus      uint8   // Result status - 0 = invalid, 1 = inactive, 2 = active, 3 = finished, 4 = disqualified, 5 = not classified, 6 = retired
          }
Total: 41

          type PacketLapData struct {
21            M_header PacketHeader // Header

820            M_lapData [20]LapData // Lap data for all cars on track
          }
Total: 841
```

## EVENT PACKET

> Size: 25 bytes

```
          type PacketEventData struct {
21            M_header PacketHeader // Header

4             M_eventStringCode [4]uint8 // Event string code, see above
        }
Total: 25
```

## PARTICIPANTS PACKET

> Size: 1082 bytes

```
          type ParticipantData struct {
1             M_aiControlled uint8      // Whether the vehicle is AI (1) or Human (0) controlled
1             M_driverId     uint8      // Driver id - see appendix
1             M_teamId       uint8      // Team id - see appendix
1             M_raceNumber   uint8      // Race number of the car
1             M_nationality  uint8      // Nationality of the driver
48            M_name         [48]string // Name of participant in UTF-8 format – null terminated, Will be truncated with … (U+2026) if too long
        }
Total:

          type PacketParticipantsData struct {
21             M_header PacketHeader // Header

1             M_numCars      uint8 // Number of cars in the data
1,060          M_participants [20]ParticipantData
        }
Total: 1082
```

## CAR SETUPS PACKET

> Size: 841 bytes

```
          type CarSetupData struct {
1             M_frontWing             uint8   // Front wing aero
1             M_rearWing              uint8   // Rear wing aero
1             M_onThrottle            uint8   // Differential adjustment on throttle (percentage)
1             M_offThrottle           uint8   // Differential adjustment off throttle (percentage)
4             M_frontCamber           float32 // Front camber angle (suspension geometry)
4             M_rearCamber            float32 // Rear camber angle (suspension geometry)
4             M_frontToe              float32 // Front toe angle (suspension geometry)
4             M_rearToe               float32 // Rear toe angle (suspension geometry)
1             M_frontSuspension       uint8   // Front suspension
1             M_rearSuspension        uint8   // Rear suspension
1             M_frontAntiRollBar      uint8   // Front anti-roll bar
1              M_rearAntiRollBar       uint8   // Front anti-roll bar
1             M_frontSuspensionHeight uint8   // Front ride height
1             M_rearSuspensionHeight  uint8   // Rear ride height
1             M_brakePressure         uint8   // Brake pressure (percentage)
1             M_brakeBias             uint8   // Brake bias (percentage)
4             M_frontTyrePressure     float32 // Front tyre pressure (PSI)
4             M_rearTyrePressure      float32 // Rear tyre pressure (PSI)
1             M_ballast               uint8   // Ballast
4             M_fuelLoad              float32 // Fuel load
        }
Total: 41

        type PacketCarSetupData struct {
21            M_header PacketHeader // Header

820            M_carSetups [20]CarSetupData
        }
Total: 841
```

## CAR TELEMETRY PACKET

> Size: 1085 bytes

```
          type CarTelemetryData struct {
2             M_speed                   uint16     // Speed of car in kilometres per hour
1             M_throttle                uint8      // Amount of throttle applied (0 to 100)
1             M_steer                   int8       // Steering (-100 (full lock left) to 100 (full lock right))
1             M_brake                   uint8      // Amount of brake applied (0 to 100)
1             M_clutch                  uint8      // Amount of clutch applied (0 to 100)
1             M_gear                    int8       // Gear selected (1-8, N=0, R=-1)
2             M_engineRPM               uint16     // Engine RPM
1             M_drs                     uint8      // 0 = off, 1 = on
1             M_revLightsPercent        uint8      // Rev lights indicator (percentage)
8             M_brakesTemperature       [4]uint16  // Brakes temperature (celsius)
8             M_tyresSurfaceTemperature [4]uint16  // Tyres surface temperature (celsius)
8             M_tyresInnerTemperature   [4]uint16  // Tyres inner temperature (celsius)
2             M_engineTemperature       uint16     // Engine temperature (celsius)
16            M_tyresPressure           [4]float32 // Tyres pressure (PSI)
        }
Total: 53

          type PacketCarTelemetryData struct {
21             M_header PacketHeader // Header

1,060          M_carTelemetryData [20]CarTelemetryData

4             M_buttonStatus uint32 // Bit flags specifying which buttons are being pressed currently - see appendices
        }
Total: 1085
```

## CAR STATUS PACKET

> Size: 1061 bytes

```
          type CarStatusData struct {
1             M_tractionControl         uint8    // 0 (off) - 2 (high)
1             M_antiLockBrakes          uint8    // 0 (off) - 1 (on)
1             M_fuelMix                 uint8    // Fuel mix - 0 = lean, 1 = standard, 2 = rich, 3 = max
1             M_frontBrakeBias          uint8    // Front brake bias (percentage)
1             M_pitLimiterStatus        uint8    // Pit limiter status - 0 = off, 1 = on
4             M_fuelInTank              float32  // Current fuel mass
4             M_fuelCapacity            float32  // Fuel capacity
2             M_maxRPM                  uint16   // Cars max RPM, point of rev limiter
2             M_idleRPM                 uint16   // Cars idle RPM
1             M_maxGears                uint8    // Maximum number of gears
1             M_drsAlloweds             uint8    // 0 = not allowed, 1 = allowed, -1 = unknown
4             M_tyresWear               [4]uint8 // Tyre wear percentage
1             M_tyreCompound            uint8    // Modern - 0 = hyper soft, 1 = ultra soft, 2 = super soft, 3 = soft, 4 = medium, 5 = hard, 6 = super hard, 7 = inter, 8 = wet, Classic - 0-6 = dry, 7-8 = wet
4             M_tyresDamage             [4]uint8 // Tyre damage (percentage)
1             M_frontLeftWingDamage     uint8    // Front left wing damage (percentage)
1             M_frontRightWingDamage    uint8    // Front right wing damage (percentage)
1             M_rearWingDamage          uint8    // Rear wing damage (percentage)
1             M_engineDamage            uint8    // Engine damage (percentage)
1             M_gearBoxDamage           uint8    // Gear box damage (percentage)
1             M_exhaustDamage           uint8    // Exhaust damage (percentage)
1             M_vehicleFiaFlags         int8     // -1 = invalid/unknown, 0 = none, 1 = green, 2 = blue, 3 = yellow, 4 = red
4             M_ersStoreEnergy          float32  // ERS energy store in Joules
1             M_ersDeployMode           uint8    // ERS deployment mode, 0 = none, 1 = low, 2 = medium, 3 = high, 4 = overtake, 5 = hotlap
4             M_ersHarvestedThisLapMGUK float32  // ERS energy harvested this lap by MGU-K
4             M_ersHarvestedThisLapMGUH float32  // ERS energy harvested this lap by MGU-H
4             M_ersDeployedThisLap      float32  // ERS energy deployed this lap
          }
Total: 52

        type PacketCarStatusData struct {
21            M_header PacketHeader // Header

1,040          M_carStatusData [20]CarStatusData
        }
```
