var popup = document.getElementById('select_from_database_popup');
// Get the button that opens the modal
var popup_select_open = document.getElementById("select_session_nav_button");
// Get the <span> element that closes the popup
var popup_close = document.getElementsByClassName("close")[0];





// Things that have to do with the select session from database popup
// When the user clicks on the button, open the modal
popup_select_open.onclick = function() {
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


function add_session_row(session_uid, session_start, session_end) {
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

  document.getElementById('popup_body').appendChild(new_div);
}


// connect to websocket
var ws = new WebSocket('ws:localhost:8080/history/ws');

// Function is called when go_websocket_server recieves a packet and sends it via the websocket
ws.onmessage = function(event) {
  var data = JSON.parse(event.data);
  console.log(data);

  var switch_number = data.M_header.M_packetId;

  switch (switch_number) {
    case 34:
      for (session_number = 0; session_number < data.Num_of_sessions; session_number++) {
        add_session_row(data.Sessions[session_number].Session_UID, data.Sessions[session_number].Session_start_time, data.Sessions[session_number].Session_end_time)
      }

      break;

  }
}
