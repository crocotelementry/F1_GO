//
// history popup
//
var history_popup = document.getElementById('select_from_database_popup');
// Get the button that opens the modal
var history_popup_open = document.getElementById("select_session_nav_button");
// Get the <span> element that closes the popup
var history_popup_close = document.getElementById("history_close");

//
// save popup
//
var save_popup = document.getElementById('save_to_database_popup');
// Get the button that opens the modal
var save_popup_open = document.getElementById("save_session_nav_button");
// Get the <span> element that closes the popup
var save_popup_close = document.getElementById("save_close");


// Get the popup_mysql_progress stuff
// saving packet
var popup_progress_title;
var popup_progress_canvas;
var popup_progress_canvas_container;
var popup_progress_ctx;
var progress_multiplier;




// Things that have to do with the select session from database popup
// When the user clicks on the button, open the modal
history_popup_open.onclick = function() {
  history_popup.style.display = "block";
}
save_popup_open.onclick = function() {
  save_popup.style.display = "block";
}
// When the user clicks on <span> (x), close the modal
history_popup_close.onclick = function() {
  history_popup.style.display = "none";
}
save_popup_close.onclick = function() {
  save_popup.style.display = "none";
}
// When the user clicks anywhere outside of the modal, close it
window.onclick = function(event) {
  if (event.target == history_popup) {
    history_popup.style.display = "none";
  }
  if (event.target == save_popup) {
    save_popup.style.display = "none";
  }
}


function add_session_row(session_uid, session_start, session_end, popup_type) {
  var new_div = document.createElement('div');
  new_div.className = 'popup_table_grid';

  var uid = document.createElement('div');
  uid.innerHTML = session_uid

  var start = document.createElement('div');
  start.innerHTML = session_start

  var end = document.createElement('div');
  end.innerHTML = session_end

  new_div.appendChild(uid)
  new_div.appendChild(start)
  new_div.appendChild(end)

  document.getElementById(popup_type).appendChild(new_div);
}


// connect to websocket
var ws = new WebSocket('ws:localhost:8080/history/ws');

// Function is called when go_websocket_server recieves a packet and sends it via the websocket
ws.onmessage = function(event) {
  var data = JSON.parse(event.data);
  console.log(data);

  var switch_number = data.M_header.M_packetId;

  switch (switch_number) {
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

    case 34:
      for (session_number = 0; session_number < data.Num_of_sessions; session_number++) {
        add_session_row(data.Sessions[session_number].Session_UID, data.Sessions[session_number].Session_start_time, data.Sessions[session_number].Session_end_time, "history_popup_body")
      }

      break;

  }
}
