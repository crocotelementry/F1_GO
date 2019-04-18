// set up the html elements
// Lap
var laps_element = document.getElementById('lap_number');
// position
var race_position_element = document.getElementById('race_position');
// last lap shit
var last_lap_time_element = document.getElementById('last_lap_time');
//  pit shit
var pit_status_element = document.getElementById('pit_status');
// sector shit`
var current_sector_element = document.getElementById('current_sector');
// drs shit
var drs_status_element = document.getElementById('drs_status');
// gear shit
var current_gear_element = document.getElementById('current_gear');
// fia flags shit
var current_fia_flags_element = document.getElementById('current_fia_flags');
// save_session_alert shit
var save_session_alert = document.getElementById('save_session_alert');
// save_session_alert_number shit
var save_session_alert_number = document.getElementById('save_session_alert_number');
// Get the save_to_database popup
var popup = document.getElementById('save_to_database_popup');
// Get the button that opens the modal
var popup_save_open = document.getElementById("save_session_nav_button");
// Get the <span> element that closes the popup
var popup_close = document.getElementsByClassName("close")[0];
// Get the popup_mysql_progress stuff
// saving packet
var popup_progress_title;
var popup_progress_canvas;
var popup_progress_canvas_container;
var popup_progress_ctx;
var progress_multiplier;


// create an array with the pit data strings
var pit_info = ['None', 'Pitting', 'In Pit Area'];

// create an array with the gear data strings
var gear_info = ['N', '1', '2', '3', '4', '5', '6', '7', '8'];
gear_info[-1] = 'R';

// create an array with the drs data strings
var drs_info = ['OFF', 'ON'];

// create an array with the fia flags data strings
var fia_info = ['None', 'Green', 'Blue', 'Yelliw', 'Red'];
fia_info[-1] = 'invalid/unknown';

// set the total laps to 0, then change when we do finally get the packet
var amount_of_laps = 0;

// get the speed_canvas container height and width
var speed_canvas_container = document.getElementById("speed_graph_container");
var rpm_canvas_container = document.getElementById("rpm_graph_container");
var gear_canvas_container = document.getElementById("gear_graph_container");
var throttle_brake_canvas_container = document.getElementById("throttle_brake_graph_container");

// Get the html5 canvas element and set the width to the chart width
var speed_canvas = document.getElementById("speed_chart_canvas");
speed_canvas.width = speed_canvas_container.offsetWidth;
speed_canvas.height = speed_canvas_container.offsetHeight;

var rpm_canvas = document.getElementById("rpm_chart_canvas");
rpm_canvas.width = rpm_canvas_container.offsetWidth;
rpm_canvas.height = rpm_canvas_container.offsetHeight;

var gear_canvas = document.getElementById("gear_chart_canvas");
gear_canvas.width = gear_canvas_container.offsetWidth;
gear_canvas.height = gear_canvas_container.offsetHeight;

var throttle_brake_canvas = document.getElementById("throttle_brake_chart_canvas");
throttle_brake_canvas.width = throttle_brake_canvas_container.offsetWidth;
throttle_brake_canvas.height = throttle_brake_canvas_container.offsetHeight;



// Get the canvas '2d' object, which can be used to draw text, lines, boxes, circles, and more - on the canvas.
// We do this since canvas doesnt actually let us draw, it is simply a container
var speed_ctx = speed_canvas.getContext("2d");
var speed_canvas_height = speed_canvas.height;
var speed_canvas_width = speed_canvas.width;

var rpm_ctx = rpm_canvas.getContext("2d");
var rpm_canvas_height = rpm_canvas.height;
var rpm_canvas_width = rpm_canvas.width;

var gear_ctx = gear_canvas.getContext("2d");
var gear_canvas_height = gear_canvas.height;
var gear_canvas_width = gear_canvas.width;

var throttle_brake_ctx = throttle_brake_canvas.getContext("2d");
var throttle_brake_canvas_height = throttle_brake_canvas.height;
var throttle_brake_canvas_width = throttle_brake_canvas.width;



// create an array that will hold our speed chart data, with the first data point being "0,0" equivilent
var speed_chart_data = [];
var rpm_chart_data = [];
var gear_chart_data = [];
var throttle_brake_chart_throttle_data = [];
var throttle_brake_chart_brake_data = [];


// create an array that will hold our speed chart data, with the first data point being "0,0" equivilent
var speed_chart_data_queue = [];
var rpm_chart_data_queue = [];
var gear_chart_data_queue = [];
var throttle_brake_chart_throttle_data_queue = [];
var throttle_brake_chart_brake_data_queue = [];

// Set a variable that will be how many points can be contained in the chart.
// For now, set this to the width of the chart so each input from our websocket will be a pixel apart
// Adjust as required
var chart_points_number = speed_canvas_width;

// Set a variable for the shift variable
var chart_shift_variable = speed_canvas_width / chart_points_number;

// translate the 0,0 point from the top left of the canvas to the bottom left of the canvas
speed_ctx.translate(0, speed_canvas_height);
speed_ctx.scale(1, -1);

rpm_ctx.translate(0, rpm_canvas_height);
rpm_ctx.scale(1, -1);

gear_ctx.translate(0, gear_canvas_height);
gear_ctx.scale(1, -1);

throttle_brake_ctx.translate(0, throttle_brake_canvas_height);
throttle_brake_ctx.scale(1, -1);

// set the chart line_widths
speed_ctx.lineWidth = "2";
rpm_ctx.lineWidth = "2";
gear_ctx.lineWidth = "2";
throttle_brake_ctx.lineWidth = "2";

// set the chart line_colors
speed_ctx.strokeStyle = "#8D8741";
rpm_ctx.strokeStyle = "#FF4136";
gear_ctx.strokeStyle = "#7FDBFF";
// throttle_brake uses two color lines, need function to draw and array to hold colors
throttle_brake_line_colors = ['#01FF70', '#FF851B'];

// Multiplier to convert speed in relation to canvas height where 0 is bottom and 350 km/h is the top
var speed_multiplier = speed_canvas_height / 350;
var rpm_multiplier = rpm_canvas_height / 12500;
var gear_multiplier = gear_canvas_height / 8;
var throttle_brake_multiplier = throttle_brake_canvas_height / 100;

// Create a variable to controll our interval. If all new data is processed then dont run the charts
var not_read_position_adder = 0;

// Things that have to do witht the save to database popup
// When the user clicks on the button, open the modal
popup_save_open.onclick = function() {
  popup.style.display = "block";
}
// When the user clicks on <span> (x), close the modal
popup_close.onclick = function() {
  popup.style.display = "none";
}
// When the user clicks anywhere outside of the modal, close it
window.onclick = function(event) {
  if (event.target == popup) {
    popup.style.display = "none";
  }
}


// connect to websocket
var ws = new WebSocket('ws://localhost:8080/ws');


// Function that shifts the chart
function chart_shift(chart_data, new_graph_points) {
  // If our chart is full with data, delete the furthest left
  // console.log(chart_points_number);
  if (chart_data.length > chart_points_number) {
    chart_data.splice(0, new_graph_points);
  };

  // shift the chart data
  for (x = 0; x < chart_data.length; x++) {
    chart_data[x][0] = chart_data[x][0] - (new_graph_points);

  };
};

// Function to clear the chart before redraw
function clear_chart(chart_ctx, chart_canvas) {
  chart_ctx.clearRect(0, 0, chart_canvas.width, chart_canvas.height);
}

// Function to draw the chart
function draw_chart_makeUp(chart_data, chart_ctx, chart_canvas, new_graph_points) {
  // First we need to clear the chart

  var previous_points = chart_data[0];

  // Begin the path of the line chart
  chart_ctx.moveTo(previous_points[0], previous_points[1]);
  chart_ctx.beginPath();


  for (x = 0; x < chart_data.length; x++) {
    chart_ctx.moveTo(previous_points[0], previous_points[1]);
    chart_ctx.lineTo(chart_data[x][0], chart_data[x][1]);
    previous_points = chart_data[x];
  };

  chart_ctx.stroke();
}

// Function to draw the chart
function draw_chart(chart_data, chart_ctx, chart_canvas, new_graph_points) {
  // First we need to clear the chart
  chart_shift(chart_data, new_graph_points);
  clear_chart(chart_ctx, chart_canvas);

  var previous_points = chart_data[0];

  // Begin the path of the line chart
  chart_ctx.moveTo(previous_points[0], previous_points[1]);
  chart_ctx.beginPath();


  for (x = 0; x < chart_data.length; x++) {
    chart_ctx.moveTo(previous_points[0], previous_points[1]);
    chart_ctx.lineTo(chart_data[x][0], chart_data[x][1]);
    previous_points = chart_data[x];
  };

  chart_ctx.stroke();
}

// Function to draw throttle_brake chart
function draw_throttle_brake_chart(throttle_data, brake_data, chart_ctx, chart_canvas, new_graph_points) {
  // First we need to clear the chart
  chart_shift(throttle_data, new_graph_points);
  chart_shift(brake_data, new_graph_points);
  clear_chart(chart_ctx, chart_canvas);

  var data_array = [throttle_data, brake_data];

  for (y = 0; y < 2; y++) {

    var previous_points = data_array[y][0];
    chart_ctx.moveTo(previous_points[0], previous_points[1]);
    chart_ctx.strokeStyle = throttle_brake_line_colors[y];
    chart_ctx.beginPath();

    for (x = 0; x < data_array[y].length; x++) {
      chart_ctx.moveTo(previous_points[0], previous_points[1]);
      chart_ctx.lineTo(data_array[y][x][0], data_array[y][x][1]);
      previous_points = data_array[y][x];
    };

    chart_ctx.stroke();

  };
}

var drawing_interval_efficiency = setInterval(draw_graphs, 18);

function draw_graphs() {
  if (not_read_position_adder > 0) {
    let num_of_new_packets = not_read_position_adder;
    not_read_position_adder = 0;
    draw_chart(speed_chart_data, speed_ctx, speed_canvas, num_of_new_packets);
    draw_chart(rpm_chart_data, rpm_ctx, rpm_canvas, num_of_new_packets);
    draw_chart(gear_chart_data, gear_ctx, gear_canvas, num_of_new_packets);
    draw_throttle_brake_chart(throttle_brake_chart_throttle_data, throttle_brake_chart_brake_data, throttle_brake_ctx, throttle_brake_canvas, num_of_new_packets);
  }
}

// Function to convert the time we are given in the UDP packets in seconds to a standard time format
function intTime_to_timeTime(time_str) {
  let step_one = time_str / 60;
  let time_min = Math.floor(step_one);

  let step_two = (step_one - time_min) * 60;
  let time_sec = Math.floor(step_two);

  let step_three = (step_two - time_sec) * 60;
  let time_mil = Math.floor(step_three);



  if (time_min < 10) {
    time_min = "0" + time_min.toString();
  } else {
    time_min = time_min.toString();
  }

  if (time_sec < 10) {
    time_sec = "0" + time_sec.toString();
  } else {
    time_sec = time_sec.toString();
  }

  if (time_mil < 10) {
    time_mil = "0" + time_mil.toString();
  } else {
    time_mil = time_mil.toString();
  }

  return time_min + ":" + time_sec + ":" + time_mil
}

function save_to_database(uid) {
  var uid_json = '{"type":"add", "uid":' + uid + '}';

  var status_body = document.getElementById(uid);

  var status_content = document.createElement('div');
  status_content.className = "popup_mysql_progress_content";


  // motion packet
  var mp_div = document.createElement('div');
  mp_div.className = "popup_progress_data_row";
  var popup_progress_data_left = document.createElement('div');
  popup_progress_data_left.id = "popup_progress_data_title";
  popup_progress_canvas_container = document.createElement('div');
  popup_progress_canvas_container.id = "popup_progress_canvas_container";
  // add title to left
  popup_progress_title = document.createElement('span');
  popup_progress_title.className = "popup_progress_title";
  // add canvas to right
  popup_progress_canvas = document.createElement('canvas');
  popup_progress_canvas.id = "popup_progress_canvas";



  popup_progress_data_left.appendChild(popup_progress_title);
  popup_progress_canvas_container.appendChild(popup_progress_canvas);
  mp_div.appendChild(popup_progress_data_left);
  mp_div.appendChild(popup_progress_canvas_container);


  status_content.appendChild(mp_div);

  status_body.appendChild(status_content);

  status_body.style.display = "block";


  ws.send(uid_json);
}

function add_session_row(session_uid, session_start, session_end) {
  var new_div = document.createElement('div');
  new_div.className = 'popup_table_grid';

  var uid = document.createElement('div');
  uid.innerHTML = session_uid

  var start = document.createElement('div');
  start.innerHTML = session_start

  var end = document.createElement('div');
  end.innerHTML = session_end

  var button_div = document.createElement('div');
  button_div.innerHTML = '<input type="button" class="save_session_button" data-uid="' + session_uid + '" value="SAVE" onclick="save_to_database(this.dataset.uid)">';

  new_div.appendChild(uid)
  new_div.appendChild(start)
  new_div.appendChild(end)
  new_div.appendChild(button_div)

  document.getElementById('popup_body').appendChild(new_div);


  var status_body = document.createElement('div');
  status_body.className = 'popup_mysql_progress';
  status_body.id = session_uid;

  document.getElementById('popup_body').appendChild(status_body);
}

// Uncomment and alter to display when the websocket is closed from the servers end
// conn.onclose = function (evt) {
//     var item = document.createElement("div");
//     item.innerHTML = "<b>Connection closed.</b>";
//     appendLog(item);
// };

// Function is called when go_websocket_server recieves a packet and sends it via the websocket
ws.onmessage = function(event) {
  var data = JSON.parse(event.data);
  // console.log(data)
  var switch_number = data.M_header.M_packetId;

  switch (switch_number) {
    // If the data inbound is the session data packet, grab the amount of total laps
    case 1:
      if (amount_of_laps == 0) {
        amount_of_laps = data.M_totalLaps;
      }
      break;

      // If the data inbound is the lap data packet
    case 2:
      // set current lap number data
      if (amount_of_laps != 0) {
        laps_element.innerHTML = JSON.stringify(data.M_lapData[data.M_header.M_playerCarIndex].M_currentLapNum, null) + " of " + JSON.stringify(amount_of_laps, null);
      }
      // Set the data for the rest lap packet informations
      race_position_element.innerHTML = data.M_lapData[data.M_header.M_playerCarIndex].M_carPosition;
      last_lap_time_element.innerHTML = intTime_to_timeTime(data.M_lapData[data.M_header.M_playerCarIndex].M_lastLapTime); // convert to time format?
      pit_status_element.innerHTML = pit_info[data.M_lapData[data.M_header.M_playerCarIndex].M_pitStatus]; // 0 = none, 1 = pitting, 2 = in pit area
      current_sector_element.innerHTML = data.M_lapData[data.M_header.M_playerCarIndex].M_sector + 1; // 0 = sector1, 1 = sector2, 2 = sector3
      break;

      // If the data inbound is the car telemetry packet
    case 6:
      // console.log(typeof (data.M_carTelemetryData))
      // Set the data for the telemetry packet information
      drs_status_element.innerHTML = drs_info[data.M_carTelemetryData[data.M_header.M_playerCarIndex].M_drs]; // 0 = off, 1 = on
      current_gear_element.innerHTML = gear_info[data.M_carTelemetryData[data.M_header.M_playerCarIndex].M_gear];

      speed_chart_data.push([speed_canvas_width + not_read_position_adder, data.M_carTelemetryData[data.M_header.M_playerCarIndex].M_speed * speed_multiplier]);
      rpm_chart_data.push([rpm_canvas_width + not_read_position_adder, data.M_carTelemetryData[data.M_header.M_playerCarIndex].M_engineRPM * rpm_multiplier]);
      gear_chart_data.push([gear_canvas_width + not_read_position_adder, data.M_carTelemetryData[data.M_header.M_playerCarIndex].M_gear * gear_multiplier]);

      throttle_brake_chart_throttle_data.push([throttle_brake_canvas_width + not_read_position_adder, data.M_carTelemetryData[data.M_header.M_playerCarIndex].M_throttle * throttle_brake_multiplier]);
      throttle_brake_chart_brake_data.push([throttle_brake_canvas_width + not_read_position_adder, data.M_carTelemetryData[data.M_header.M_playerCarIndex].M_brake * throttle_brake_multiplier]);
      //
      not_read_position_adder += 1;
      break;


      // If the data inbound is the car status packet
    case 7:
      // Set the data for the status packet information
      // console.log(switch_number, data)
      // console.log(typeof (data.M_carStatusData));
      current_fia_flags_element.innerHTML = fia_info[data.M_carStatusData[data.M_header.M_playerCarIndex].M_vehicleFiaFlags];
      break; // -1 = invalid/unknown, 0 = none, 1 = green, 2 = blue, 3 = yellow, 4 = red

    case 30:
      //
      // console.log("Session just finished, sending redis captured info");
      // console.log("Num of session", data.Num_of_sessions);
      // console.log("Session 1 data:", data.Sessions[0].Session_UID, data.Sessions[0].Session_start_time, data.Sessions[0].Session_end_time);
      save_session_alert_number.innerHTML = data.Num_of_sessions;
      save_session_alert.classList.toggle('show');
      save_session_alert.classList.toggle('hide');

      for (z = 0; z < data.Num_of_sessions; z++) {
        add_session_row(data.Sessions[z].Session_UID, data.Sessions[z].Session_start_time, data.Sessions[z].Session_end_time)
      }

      break;

      // Case 31, if data inbound is for Save_to_database_status
    case 31:
      //
      // console.log("MYSQL PROGRESS DATA:", data);

      switch (data.Status) {
        case "initial":
          // popup_progress_canvas
          popup_progress_total = data.Total_packets;
          // popup_progress_canvas = document.getElementById("popup_progress_canvas");
          popup_progress_canvas.width = popup_progress_canvas_container.offsetWidth;
          popup_progress_canvas.height = popup_progress_canvas_container.offsetHeight;
          // Get the canvas '2d' object, which can be used to draw text, lines, boxes, circles, and more - on the canvas.
          // We do this since canvas doesnt actually let us draw, it is simply a container
          popup_progress_ctx = popup_progress_canvas.getContext("2d");
          var popup_progress_canvas_height = popup_progress_canvas.height;
          var popup_progress_canvas_width = popup_progress_canvas.width;
          popup_progress_ctx.fillStyle = "#DDDDDD";
          // Multiplier to convert speed in relation to canvas height where 0 is bottom and 350 km/h is the top
          progress_multiplier = popup_progress_canvas_width / popup_progress_total;
          popup_progress_ctx.fillRect(0, 0, progress_multiplier * data.Total_current, popup_progress_canvas.height);
          break;

        case "Saving":
          switch (data.Current_packet) {
            case 0:
              if (popup_progress_title.innerHTML != "Saving: Motion Packet") {
                popup_progress_title.innerHTML = "Motion Packet";
              }
              popup_progress_ctx.fillRect(0, 0, progress_multiplier * data.Total_current, popup_progress_canvas.height);
              break;
            case 1:
              if (popup_progress_title.innerHTML != "Saving: Session Packet") {
                popup_progress_title.innerHTML = "Session Packet";
              }
              popup_progress_ctx.fillRect(0, 0, progress_multiplier * data.Total_current, popup_progress_canvas.height);
              break;
            case 2:
              if (popup_progress_title.innerHTML != "Saving: Lap Data Packet") {
                popup_progress_title.innerHTML = "Lap Data Packet";
              }
              popup_progress_ctx.fillRect(0, 0, progress_multiplier * data.Total_current, popup_progress_canvas.height);
              break;
            case 3:
              if (popup_progress_title.innerHTML != "Saving: Event Packet") {
                popup_progress_title.innerHTML = "Event Packet";
              }
              popup_progress_ctx.fillRect(0, 0, progress_multiplier * data.Total_current, popup_progress_canvas.height);
              break;
            case 4:
              if (popup_progress_title.innerHTML != "Saving: Participants Packet") {
                popup_progress_title.innerHTML = "Participants Packet";
              }
              popup_progress_ctx.fillRect(0, 0, progress_multiplier * data.Total_current, popup_progress_canvas.height);
              break;
            case 5:
              if (popup_progress_title.innerHTML != "Saving: Car Setups Packet") {
                popup_progress_title.innerHTML = "Car Setups Packet";
              }
              popup_progress_ctx.fillRect(0, 0, progress_multiplier * data.Total_current, popup_progress_canvas.height);
              break;
            case 6:
              if (popup_progress_title.innerHTML != "Saving: Car Telemety Packet") {
                popup_progress_title.innerHTML = "Car Telemetry Packet";
              }
              popup_progress_ctx.fillRect(0, 0, progress_multiplier * data.Total_current, popup_progress_canvas.height);
              break;
            case 7:
              if (popup_progress_title.innerHTML != "Saving: Car Status Packet") {
                popup_progress_title.innerHTML = "Car Status Packet";
              }
              popup_progress_ctx.fillRect(0, 0, progress_multiplier * data.Total_current, popup_progress_canvas.height);
              break;
          }
          break;

        case "done":
          popup_progress_title.innerHTML = "Completed"
          console.log("Finished adding redis data for session uid to mysql database:", data.UID);
          //
          //
          // TAKE AWAY SAVE BUTTON AND ADD REMOVE BUTTON
          //
          //
          break;

        default:
          popup_progress_title.innerHTML = "Error!"
          console.log("error with popup_progress data.Status switch statement");
          break;
      }
      break;

    case 32:

      for (i = 0; i < (data.RaceSpeed_data).length; i++) {
        speed_chart_data.push([speed_canvas_width + not_read_position_adder, data.RaceSpeed_data[((data.RaceSpeed_data).length - 1) - i] * speed_multiplier]);
        rpm_chart_data.push([rpm_canvas_width + not_read_position_adder, data.EngineRevs_data[((data.EngineRevs_data).length - 1) - i] * rpm_multiplier]);
        gear_chart_data.push([gear_canvas_width + not_read_position_adder, data.GearChanges_data[((data.GearChanges_data).length - 1) - i] * gear_multiplier]);
        throttle_brake_chart_throttle_data.push([throttle_brake_canvas_width + not_read_position_adder, data.ThrottleApplication_data[((data.ThrottleApplication_data).length - 1) - i] * throttle_brake_multiplier]);
        throttle_brake_chart_brake_data.push([throttle_brake_canvas_width + not_read_position_adder, data.BrakeApplication_data[((data.BrakeApplication_data).length - 1) - i] * throttle_brake_multiplier]);

        not_read_position_adder += 1;
      }

      // console.log(switch_number, "catchUP!!!");
      break;

    default:
      console.log(switch_number, "Invalid packeted id sent over websocket!\n", data);
      break;
  }
}
