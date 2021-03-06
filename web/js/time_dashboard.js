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


var new_time_entry_sector_one = '00:00:00';
var new_time_entry_sector_two = '00:00:00';
var new_time_entry_sector_three = '00:00:00';
var new_time_entry_lap_time = '00:00:00';

// var new_time_entry_html_part_one    = '<div class="time_grid_entry new_lap_entry" id="';
// var new_time_entry_html_part_two    = '"><div class="time_grid_data"><div class="lap_number"><span class="lap_number_text">';
// var new_time_entry_html_part_three  = '</span></div><div class="sector_one"><span class="sector_one_text" id="';
// var new_time_entry_html_part_four   = '_sector_one_text">';
// var new_time_entry_html_part_five   = '</span></div><div class="sector_two"><span class="sector_two_text" id="';
// var new_time_entry_html_part_six    = '_sector_two_text">';
// var new_time_entry_html_part_seven  = '</span></div><div class="sector_three"><span class="sector_three_text" id="';
// var new_time_entry_html_part_eight  = '_sector_three_text">';
// var new_time_entry_html_part_nine   = '</span></div><div class="lap_time"><span class="lap_time_text" id="';
// var new_time_entry_html_part_ten    = '_lap_time_text">';
// var new_time_entry_html_part_eleven  = '</span></div></div></div>';

var current_lap = 0;

var one = current_lap.toString() + '_sector_one_text';
var two = current_lap.toString() + '_sector_two_text'
var three = current_lap.toString() + '_sector_three_text'
var four = current_lap.toString() + '_lap_time_text'

var current_time_entry_sector_one;
var current_time_entry_sector_two;
var current_time_entry_sector_three;
var current_time_entry_lap_time;

// zero if hasnt gone to pit, 1 if car has
var has_pitted = 0;

// Fastest lap
var fastest_lap_time = 0;
var fastest_lap_number = 0;

// Slowest lap
var slowest_lap_time = 0;
var slowest_lap_number = 0;

// Current lap time
var current_lap_time;


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


// <div class="time_grid_entry new_lap_entry" id="">
//   <div class="time_grid_data">
//     <div class="lap_number">
//       <span class="lap_number_text"></span>
//     </div>
//     <div class="sector_one">
//       <span class="sector_one_text" id="_sector_one_text"></span>
//     </div>
//     <div class="sector_two">
//       <span class="sector_two_text" id="_sector_two_text"></span>
//     </div>
//     <div class="sector_three">
//       <span class="sector_three_text" id="_sector_three_text"></span>\
//     </div>
//     <div class="lap_time">
//       <span class="lap_time_text" id="_lap_time_text"></span>
//     </div>
//   </div>
// </div>

function add_time_row(lap_number) {
  var new_lap_row = document.createElement('div');
  new_lap_row.className = 'time_grid_entry new_lap_entry';
  new_lap_row.id = "lap_" + lap_number;

  var new_lap_data = document.createElement('div');
  new_lap_data.className = 'time_grid_data';

  var new_lap_number = document.createElement('div');
  new_lap_number.className = 'lap_number';

  var new_lap_number_text = document.createElement('span');
  new_lap_number_text.className = 'lap_number_text';
  new_lap_number_text.innerHTML = lap_number;


  var new_sector_one = document.createElement('div');
  new_sector_one.className = 'sector_one';

  var new_sector_one_text = document.createElement('span');
  new_sector_one_text.className = 'sector_one_text';
  new_sector_one_text.id = lap_number + '_sector_one_text';
  new_sector_one_text.innerHTML = new_time_entry_sector_one;


  var new_sector_two = document.createElement('div');
  new_sector_two.className = 'sector_two';

  var new_sector_two_text = document.createElement('span');
  new_sector_two_text.className = 'sector_two_text';
  new_sector_two_text.id = lap_number + '_sector_two_text';
  new_sector_two_text.innerHTML = new_time_entry_sector_two;


  var new_sector_three = document.createElement('div');
  new_sector_three.className = 'sector_three';

  var new_sector_three_text = document.createElement('span');
  new_sector_three_text.className = 'sector_three_text';
  new_sector_three_text.id = lap_number + '_sector_three_text';
  new_sector_three_text.innerHTML = new_time_entry_sector_three;


  var new_lap_time = document.createElement('div');
  new_lap_time.className = 'lap_time';

  var new_lap_time_text = document.createElement('span');
  new_lap_time_text.className = 'lap_time_text';
  new_lap_time_text.id = lap_number + '_lap_time_text';
  new_lap_time_text.innerHTML = new_time_entry_lap_time;


  // Append from the most nested element upwards. Starting with adding spans to divs and so forth
  new_lap_number.appendChild(new_lap_number_text);
  new_sector_one.appendChild(new_sector_one_text);
  new_sector_two.appendChild(new_sector_two_text);
  new_sector_three.appendChild(new_sector_three_text);
  new_lap_time.appendChild(new_lap_time_text);

  // Now append upwards one step

  new_lap_data.appendChild(new_lap_number);
  new_lap_data.appendChild(new_sector_one);
  new_lap_data.appendChild(new_sector_two);
  new_lap_data.appendChild(new_sector_three);
  new_lap_data.appendChild(new_lap_time);

  // Now append to the top
  new_lap_row.appendChild(new_lap_data);

  // Finally append new lap row to the time chart grid
  document.getElementById('time_chart_grid').appendChild(new_lap_row);
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


// connect to websocket
var ws = new WebSocket('ws:localhost:8080/time/ws');

// Function is called when go_websocket_server recieves a packet and sends it via the websocket
ws.onmessage = function(event) {
  var data = JSON.parse(event.data);

  switch (data.M_header.M_packetId) {
    // If the data inbound is the lap data packet, grab the amount of total laps
    case 2:
      if (data.M_lapData[data.M_header.M_playerCarIndex].M_currentLapNum > current_lap) {

        // Check if lap is fastest lap and at same time check for slowest lap
        //
        // Check to make sure we have actually finished a full lap first. Not just jumping to a lap mid race like with the pcap or
        // starting this on a race that is currently in progress. Without this we would set fastest lap time to 0 and would never
        // actually find the fastest lap
        if (current_lap != 0) {
          // If we have already had a full lap and our fastest lap time and fastest lap is still 0 meaning we havent had a fastest lap
          // yet becuase this would be our first real compleated lap. Set this first lap to our fastest lap.
          if (fastest_lap_time == 0) {
            fastest_lap_time = current_lap_time;
            fastest_lap_number = current_lap;
            // toggle the color for the row of the fastest lap
            document.getElementById('lap_' + fastest_lap_number).classList.toggle('fastest_lap');

            // Also set this as our slowest lap time since this would be our first completed lap, meaning it is not just our fastest lap
            // time, it is also our slowest lap time
            slowest_lap_time = current_lap_time;
            slowest_lap_number = current_lap;
          }
          // If our fastest lap time is not zero, meaning we have already compleated and set our first lap as our fastest. Now compare
          // current compleated lap with our fastest and see if we now have a new fastest lap!
          // Also, since we have now compleated two laps, we can now set the slowest of the two as our slowest lap. Eliminated the need for
          // more confusing logic statements, just check and do stuff here for it.
          else {
            // If current lap is faster than our recorded fastest lap time, we have found a new fastest lap!
            if (current_lap_time < fastest_lap_time) {
              // Toggle the color off the previous fastest lap
              document.getElementById('lap_' + fastest_lap_number).classList.toggle('fastest_lap');
              fastest_lap_time = current_lap_time;
              fastest_lap_number = current_lap;
              // Now after setting the new fastest lap number, toggle that lap numbers color on
              document.getElementById('lap_' + fastest_lap_number).classList.toggle('fastest_lap');
            }

            // Check if our past full lap is now the slowest
            if (current_lap_time > slowest_lap_time) {
              slowest_lap_time = current_lap_time;
              slowest_lap_number = current_lap;
              document.getElementById('lap_' + slowest_lap_number).classList.toggle('slowest_lap');
            }
            // In the wierd case in which the first lap is actually our slowest lap. We need to toggle its color for slowest lap since,
            // we arent going to toggle both the fastest and the slowest on the same lap after only one lap. Fastest takes preseident because it
            // truly is the fastest yet. So now toggle the color on the first lap.
            else {
              if (!document.getElementById('lap_' + slowest_lap_number).classList.contains('slowest_lap')) {
                document.getElementById('lap_' + slowest_lap_number).classList.toggle('slowest_lap');
              }
            }
          }
        }

        // Check if we pitted this lap
        if (has_pitted != 0) {
          document.getElementById('lap_' + current_lap).classList.toggle('car_pitted');
          // Set has pited to zero since we have already checked if we have, and now we are starting a new lap
          has_pitted = 0;
        }

        // // Now set the current lap to the new lap we are starting
        current_lap = data.M_lapData[data.M_header.M_playerCarIndex].M_currentLapNum;
        // // Create the string that is the html for the new lap to be added to our time chart
        // var new_time_entry = new_time_entry_html_part_one + 'lap_' + current_lap.toString() + new_time_entry_html_part_two + current_lap.toString() + new_time_entry_html_part_three + current_lap.toString() + new_time_entry_html_part_four + new_time_entry_sector_one + new_time_entry_html_part_five + current_lap.toString() + new_time_entry_html_part_six + new_time_entry_sector_two + new_time_entry_html_part_seven + current_lap.toString() + new_time_entry_html_part_eight + new_time_entry_sector_three + new_time_entry_html_part_nine + current_lap.toString() + new_time_entry_html_part_ten + new_time_entry_lap_time + new_time_entry_html_part_eleven;
        // // Add the new laps html row into our time chart
        // document.getElementById('time_chart_grid').innerHTML += new_time_entry;
        add_time_row(current_lap.toString())

        // Get the elements for the new laps html elements representing the time values
        current_time_entry_sector_one = document.getElementById(current_lap.toString() + '_sector_one_text');
        current_time_entry_sector_two = document.getElementById(current_lap.toString() + '_sector_two_text');
        current_time_entry_sector_three = document.getElementById(current_lap.toString() + '_sector_three_text');
        current_time_entry_lap_time = document.getElementById(current_lap.toString() + '_lap_time_text');
      }

      // Set the html elements for the current lap to the times sent over from our udp packets
      current_time_entry_lap_time.innerHTML = intTime_to_timeTime(data.M_lapData[data.M_header.M_playerCarIndex].M_currentLapTime);
      current_lap_time = data.M_lapData[data.M_header.M_playerCarIndex].M_currentLapTime;
      // Since there is no sector three time udp packet data, we can check what sector we are in, and if we are in sector three then
      // we subtract the two sector times from our overall current lap time to get the time for sector three
      current_time_entry_sector_one.innerHTML = intTime_to_timeTime(data.M_lapData[data.M_header.M_playerCarIndex].M_sector1Time);
      current_time_entry_sector_two.innerHTML = intTime_to_timeTime(data.M_lapData[data.M_header.M_playerCarIndex].M_sector2Time);
      if (data.M_lapData[data.M_header.M_playerCarIndex].M_sector == 2) {
        current_time_entry_sector_three.innerHTML = intTime_to_timeTime(data.M_lapData[data.M_header.M_playerCarIndex].M_currentLapTime - data.M_lapData[data.M_header.M_playerCarIndex].M_sector2Time - data.M_lapData[data.M_header.M_playerCarIndex].M_sector1Time);
      }

      // Only check if we have pitted when we havent. Save the computer from checking if we have the status and setting the piting to 1 if we already have
      if (has_pitted == 0) {
        if (data.M_lapData[data.M_header.M_playerCarIndex].M_pitStatus == 1 || data.M_lapData[data.M_header.M_playerCarIndex].M_pitStatus == 2) {
          has_pitted = 1;
        }
      }


      // End of case thingy stuff stuff
      break;

      // Case 30, if data inbound is for finished sessions waiting to be saved to database
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

      // If we get time catchup data
    case 33:
      if ((data.Lap_num).length > 0) {
        console.log("recieved catchup data")
        console.log(data)

        catchup_fast_lap_time = data.Lap_time[0];
        catchup_fast_lap_number = data.Lap_num[0];
        catchup_slow_lap_time = data.Lap_time[0];
        catchup_slow_lap_number = data.Lap_num[0];

        for (i = 0; i < (data.Lap_num).length; i++) {
          current_catchup_lap = data.Lap_num[i];
          add_time_row(current_catchup_lap.toString());

          current_time_entry_sector_one = document.getElementById(current_catchup_lap.toString() + '_sector_one_text');
          current_time_entry_sector_two = document.getElementById(current_catchup_lap.toString() + '_sector_two_text');
          current_time_entry_sector_three = document.getElementById(current_catchup_lap.toString() + '_sector_three_text');
          current_time_entry_lap_time = document.getElementById(current_catchup_lap.toString() + '_lap_time_text');

          current_time_entry_lap_time.innerHTML = intTime_to_timeTime(data.Lap_time[i]);
          current_time_entry_sector_one.innerHTML = intTime_to_timeTime(data.Sector1Time[i]);
          current_time_entry_sector_two.innerHTML = intTime_to_timeTime(data.Sector2Time[i]);
          current_time_entry_sector_three.innerHTML = intTime_to_timeTime(data.Sector3Time[i]);

          if (data.PitStatus[i] == 1 || data.PitStatus[i] == 2) {
            document.getElementById('lap_' + current_catchup_lap).classList.toggle('car_pitted');
          }

          if (data.Lap_time[i] > catchup_slow_lap_time) {
            catchup_slow_lap_time = data.Lap_time[i];
            catchup_slow_lap_number = data.Lap_num[i];
          }

          if (data.Lap_time[i] < catchup_fast_lap_time) {
            catchup_fast_lap_time = data.Lap_time[i];
            catchup_fast_lap_number = data.Lap_num[i];
          }

        }


        fastest_lap_time = catchup_fast_lap_time;
        fastest_lap_number = catchup_fast_lap_number;

        console.log("catchup_fast_lap_number:", catchup_fast_lap_number);
        // Now after setting the new fastest lap number, toggle that lap numbers color on
        document.getElementById('lap_' + catchup_fast_lap_number).classList.toggle('fastest_lap');

        slowest_lap_time = catchup_slow_lap_time;
        slowest_lap_number = catchup_slow_lap_number;
        document.getElementById('lap_' + catchup_slow_lap_number).classList.toggle('slowest_lap');

        current_lap = i;
      }

      break;



  }
}
