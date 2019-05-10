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


// playback options
var playback_start_button = document.getElementById('playback_start');
var playback_stop_button = document.getElementById('playback_stop');
var playback_speed_normal_button = document.getElementById('playback_speed_normal');
var playback_speed_half_button = document.getElementById('playback_speed_half');
var playback_speed_quarter_button = document.getElementById('playback_speed_quarter');

var previous_popup_speed = 0;
var playback_speed_list = [playback_speed_normal_button, playback_speed_half_button, playback_speed_quarter_button];

// Get our html elements for our history data
// times
var total_time = document.getElementById('total_time');
var current_lap = document.getElementById('current_lap');
var last_lap = document.getElementById('last_lap');
var best_lap = document.getElementById('best_lap');
var penalties = document.getElementById('penalties');
var current_sector = document.getElementById('current_sector');
var sector_1 = document.getElementById('sector_1');
var sector_2 = document.getElementById('sector_2');
var sector_3 = document.getElementById('sector_3');
// standings
var standing_color_1 = document.getElementById('standing_color_1');
var standing_color_2 = document.getElementById('standing_color_2');
var standing_color_3 = document.getElementById('standing_color_3');
var standing_color_4 = document.getElementById('standing_color_4');
var standing_color_5 = document.getElementById('standing_color_5');
var standing_color_6 = document.getElementById('standing_color_6');
var standing_color_7 = document.getElementById('standing_color_7');
var standing_color_8 = document.getElementById('standing_color_8');
var standing_color_9 = document.getElementById('standing_color_9');
var standing_color_10 = document.getElementById('standing_color_10');
var standing_color_11 = document.getElementById('standing_color_11');
var standing_color_12 = document.getElementById('standing_color_12');
var standing_color_13 = document.getElementById('standing_color_13');
var standing_color_14 = document.getElementById('standing_color_14');
var standing_color_15 = document.getElementById('standing_color_15');
var standing_color_16 = document.getElementById('standing_color_16');
var standing_color_17 = document.getElementById('standing_color_17');
var standing_color_18 = document.getElementById('standing_color_18');
var standing_color_19 = document.getElementById('standing_color_19');
var standing_color_20 = document.getElementById('standing_color_20');

var standing_name_1 = document.getElementById('standing_name_1');
var standing_name_2 = document.getElementById('standing_name_2');
var standing_name_3 = document.getElementById('standing_name_3');
var standing_name_4 = document.getElementById('standing_name_4');
var standing_name_5 = document.getElementById('standing_name_5');
var standing_name_6 = document.getElementById('standing_name_6');
var standing_name_7 = document.getElementById('standing_name_7');
var standing_name_8 = document.getElementById('standing_name_8');
var standing_name_9 = document.getElementById('standing_name_9');
var standing_name_10 = document.getElementById('standing_name_10');
var standing_name_11 = document.getElementById('standing_name_11');
var standing_name_12 = document.getElementById('standing_name_12');
var standing_name_13 = document.getElementById('standing_name_13');
var standing_name_14 = document.getElementById('standing_name_14');
var standing_name_15 = document.getElementById('standing_name_15');
var standing_name_16 = document.getElementById('standing_name_16');
var standing_name_17 = document.getElementById('standing_name_17');
var standing_name_18 = document.getElementById('standing_name_18');
var standing_name_19 = document.getElementById('standing_name_19');
var standing_name_20 = document.getElementById('standing_name_20');

var standing_time_1 = document.getElementById('standing_time_1');
var standing_time_2 = document.getElementById('standing_time_2');
var standing_time_3 = document.getElementById('standing_time_3');
var standing_time_4 = document.getElementById('standing_time_4');
var standing_time_5 = document.getElementById('standing_time_5');
var standing_time_6 = document.getElementById('standing_time_6');
var standing_time_7 = document.getElementById('standing_time_7');
var standing_time_8 = document.getElementById('standing_time_8');
var standing_time_9 = document.getElementById('standing_time_9');
var standing_time_10 = document.getElementById('standing_time_10');
var standing_time_11 = document.getElementById('standing_time_11');
var standing_time_12 = document.getElementById('standing_time_12');
var standing_time_13 = document.getElementById('standing_time_13');
var standing_time_14 = document.getElementById('standing_time_14');
var standing_time_15 = document.getElementById('standing_time_15');
var standing_time_16 = document.getElementById('standing_time_16');
var standing_time_17 = document.getElementById('standing_time_17');
var standing_time_18 = document.getElementById('standing_time_18');
var standing_time_19 = document.getElementById('standing_time_19');
var standing_time_20 = document.getElementById('standing_time_20');


var standing_color_list = [standing_color_1, standing_color_2, standing_color_3, standing_color_4, standing_color_5, standing_color_6, standing_color_7, standing_color_8, standing_color_9, standing_color_10, standing_color_11, standing_color_12, standing_color_13, standing_color_14, standing_color_15, standing_color_16, standing_color_17, standing_color_18, standing_color_19, standing_color_20]

var standing_name_list = [standing_name_1, standing_name_2, standing_name_3, standing_name_4, standing_name_5, standing_name_6, standing_name_7, standing_name_8, standing_name_9, standing_name_10, standing_name_11, standing_name_12, standing_name_13, standing_name_14, standing_name_15, standing_name_16, standing_name_17, standing_name_18, standing_name_19, standing_name_20];

var standing_time_list = [standing_time_1, standing_time_2, standing_time_3, standing_time_4, standing_time_5, standing_time_6, standing_time_7, standing_time_8, standing_time_9, standing_time_10, standing_time_11, standing_time_12, standing_time_13, standing_time_14, standing_time_15, standing_time_16, standing_time_17, standing_time_18, standing_time_19, standing_time_20];

var participantData_dict = [];

// tires
var fl_tyre_pressure_data = document.getElementById('fl_tyre_pressure_data');
var fl_tyre_wear_data = document.getElementById('fl_tyre_wear_data');
var fl_tyre_temp_data = document.getElementById('fl_tyre_temp_data');
var fl_suspension_position_data = document.getElementById('fl_suspension_position_data');
var fl_tyre_damage_data = document.getElementById('fl_tyre_damage_data');
var fl_break_temp_data = document.getElementById('fl_break_temp_data');

var bl_suspension_position_data = document.getElementById('bl_suspension_position_data');
var bl_tyre_damage_data = document.getElementById('bl_tyre_damage_data');
var bl_break_temp_data = document.getElementById('bl_break_temp_data');
var bl_tyre_pressure_data = document.getElementById('bl_tyre_pressure_data');
var bl_tyre_wear_data = document.getElementById('bl_tyre_wear_data');
var bl_tyre_temp_data = document.getElementById('bl_tyre_temp_data');

var fr_tyre_temp_data = document.getElementById('fr_tyre_temp_data');
var fr_tyre_wear_data = document.getElementById('fr_tyre_wear_data');
var fr_tyre_pressure_data = document.getElementById('fr_tyre_pressure_data');
var fr_break_temp_data = document.getElementById('fr_break_temp_data');
var fr_tyre_damage_data = document.getElementById('fr_tyre_damage_data');
var fr_suspension_position_data = document.getElementById('fr_suspension_position_data');

var br_break_temp_data = document.getElementById('br_break_temp_data');
var br_tyre_damage_data = document.getElementById('br_tyre_damage_data');
var br_suspension_position_data = document.getElementById('br_suspension_position_data');
var br_tyre_temp_data = document.getElementById('br_tyre_temp_data');
var br_tyre_wear_data = document.getElementById('br_tyre_wear_data');
var br_tyre_pressure_data = document.getElementById('br_tyre_pressure_data');

// Get the popup_mysql_progress stuff
// saving packet
var popup_progress_title;
var popup_progress_canvas;
var popup_progress_canvas_container;
var popup_progress_ctx;
var progress_multiplier;

// get the sgr_canvas container height and width
var sgr_canvas_container = document.getElementById("sgr_graph_container");

// Get the html5 canvas element and set the width to the chart width
var sgr_canvas = document.getElementById("sgr_graph_canvas");
sgr_canvas.width = sgr_canvas_container.offsetWidth;
sgr_canvas.height = sgr_canvas_container.offsetHeight;

var sgr_playback_line_canvas = document.getElementById("sgr_graph_playback_line_canvas");
sgr_playback_line_canvas.width = sgr_canvas.width;
sgr_playback_line_canvas.height = sgr_canvas.height;

// Get the canvas '2d' object, which can be used to draw text, lines, boxes, circles, and more - on the canvas.
// We do this since canvas doesnt actually let us draw, it is simply a container
var sgr_ctx = sgr_canvas.getContext("2d");
var sgr_canvas_height = sgr_canvas.height;
var sgr_canvas_width = sgr_canvas.width;

var sgr_playback_line_ctx = sgr_playback_line_canvas.getContext("2d");
var sgr_playback_line_canvas_height = sgr_playback_line_canvas.height;
var sgr_playback_line_canvas_width = sgr_playback_line_canvas.width;

// create an array that will hold our sgr chart data, with the first data point being "0,0" equivilent
var sgr_chart_data = [];
var speed_data_array = [];
var gear_data_array = [];
var rpm_data_array = [];

// create an array that will hold our sgr chart data queue, with the first data point being "0,0" equivilent
// var sgr_chart_data_queue = [];

// Set a variable that will be how many points can be contained in the chart.
// For now, set this to the width of the chart so each input from our websocket will be a pixel apart
// Adjust as required
var chart_points_number = sgr_canvas_width;

// translate the 0,0 point from the top left of the canvas to the bottom left of the canvas
sgr_ctx.translate(0, sgr_canvas_height);
sgr_ctx.scale(1, -1);

sgr_playback_line_ctx.translate(0, sgr_playback_line_canvas_height);
sgr_playback_line_ctx.scale(1, -1);

// set the chart line_widths
sgr_ctx.lineWidth = "2";

sgr_playback_line_ctx.lineWidth = "2";

// sgr uses three color lines, need function to draw and array to hold colors
sgr_line_colors = ["#8D8741", "#7FDBFF", "#FF4136"];

sgr_playback_line_color = "white";

// Multiplier to convert graph items to a form where the min value is the bottom of the canvas and the
// max is the top of the canvas.
var speed_multiplier = sgr_canvas_height / 350;
var rpm_multiplier = sgr_canvas_height / 12500;
var gear_multiplier = sgr_canvas_height / 8;


// global variable for the frame we are currently playing back from
var playback_id = 0;
var frame_number = 0;
var selected_play_back_speed = 1;

// Arrays that will hold our history data
var motionData_array = [];
var sessionData_array = [];
var lapData_array = [];
var participantData_array = [];
var telemetryData_array = [];
var statusData_array = [];
var standings_standingsData_array = [];
var standings_lapDataTimes_array = [];


// Holds data for track id and total laps
var totalLaps
var trackId

// placeholders for our data Arrays
var packet_number_motionData_array = 0;
var packet_number_lapData_array = 0;
var packet_number_telemetryData_array = 0;
var packet_number_statusData_array = 0;
var packet_number_standings_standingsData_array = 0;
var packet_number_standings_lapDataTimes_array = 0;


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



var bottom_section_grid_element_hover_over = function() {
  var svg_to_color = this.dataset.hover_object;
  document.getElementById(svg_to_color).setAttribute("fill", "#8D8741");;
}

var bottom_section_grid_element_hover_out = function() {
  var svg_to_color = this.dataset.hover_object;
  document.getElementById(svg_to_color).setAttribute("fill", "none");
}

var bottom_section_grid_elements = document.getElementsByClassName("bottom_section_grid_element");

for (var i = 0; i < bottom_section_grid_elements.length; i++) {
  bottom_section_grid_elements[i].addEventListener("mouseover", bottom_section_grid_element_hover_over);
  bottom_section_grid_elements[i].addEventListener("mouseout", bottom_section_grid_element_hover_out);
}



playback_start_button.onclick = function() {
  playback_start_button.classList.add('selected_playback_option');
  playback_stop_button.classList.remove('selected_playback_option');

  if (playback_id == 0) {
    console.log("playback_started");
    start_playback(frame_number, selected_play_back_speed);
  }
}

playback_stop_button.onclick = function() {
  playback_stop_button.classList.add('selected_playback_option');
  playback_start_button.classList.remove('selected_playback_option');
  stop_playback();
}

playback_speed_normal_button.onclick = function() {
  playback_speed_normal_button.classList.add('selected_playback_option');
  // playback_speed_list[previous_popup_speed].classList.remove('selected_playback_option');

  if (selected_play_back_speed != 1) {
    playback_speed_list[previous_popup_speed].classList.remove('selected_playback_option');
  }

  if (playback_id != 0 && selected_play_back_speed != 1) {
    stop_playback();
    selected_play_back_speed = 1;
    start_playback(frame_number, selected_play_back_speed);
  } else {
    selected_play_back_speed = 1;
  }

  previous_popup_speed = 0;
}

playback_speed_half_button.onclick = function() {
  playback_speed_half_button.classList.add('selected_playback_option');
  // playback_speed_list[previous_popup_speed].classList.remove('selected_playback_option');

  if (selected_play_back_speed != 2) {
    playback_speed_list[previous_popup_speed].classList.remove('selected_playback_option');
  }

  if (playback_id != 0 && selected_play_back_speed != 2) {
    stop_playback();
    selected_play_back_speed = 2;
    start_playback(frame_number, selected_play_back_speed);
  } else {
    selected_play_back_speed = 2;
  }

  previous_popup_speed = 1;
}

playback_speed_quarter_button.onclick = function() {
  playback_speed_quarter_button.classList.add('selected_playback_option');
  // playback_speed_list[previous_popup_speed].classList.remove('selected_playback_option');

  if (selected_play_back_speed != 4) {
    playback_speed_list[previous_popup_speed].classList.remove('selected_playback_option');
  }

  if (playback_id != 0 && selected_play_back_speed != 4) {
    stop_playback();
    selected_play_back_speed = 4;
    start_playback(frame_number, selected_play_back_speed);
  } else {
    selected_play_back_speed = 4;
  }

  previous_popup_speed = 2;
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

function get_session_from_mysql(session_uid) {
  var uid_json = '{"type":"get_history", "uid":' + session_uid + '}';
  ws.send(uid_json);
}


function add_session_row(session_uid, session_start, session_end, popup_type) {
  var new_div = document.createElement('div');
  new_div.className = 'popup_table_grid';


  switch (popup_type) {
    case "save_popup_body":
      new_div.onclick = function() {
        save_to_database(session_uid);
      }
      break;

    case "history_popup_body":
      new_div.onclick = function() {
        get_session_from_mysql(session_uid);
      }
      break;
  }


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



function add_lap_to_lap_selector(lap_num) {
  var new_div = document.createElement('div');
  new_div.className = 'lap_selection_element_class';

  var new_span = document.createElement('span');
  new_span.className = 'select_lap_nav_button';
  new_span.innerHTML = "Lap " + lap_num;

  new_div.appendChild(new_span);

  document.getElementById('lap_selection_left').appendChild(new_div);
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

// Function to clear the chart before redraw
function clear_chart(chart_ctx, chart_canvas) {
  chart_ctx.clearRect(0, 0, chart_canvas.width, chart_canvas.height);
}

// Function to draw throttle_brake chart
function draw_sgr_chart(speed_data, gear_data, rev_data, chart_ctx, chart_canvas) {
  // First we need to clear the chart
  // chart_shift(throttle_data, new_graph_points);
  // chart_shift(brake_data, new_graph_points);
  clear_chart(chart_ctx, chart_canvas);

  var data_array = [speed_data, gear_data, rev_data];

  for (data_array_position = 0; data_array_position < 3; data_array_position++) {

    var previous_points = data_array[data_array_position][0];
    chart_ctx.moveTo(previous_points[0], previous_points[1]);
    chart_ctx.strokeStyle = sgr_line_colors[data_array_position];
    chart_ctx.beginPath();

    for (x = 0; x < data_array[data_array_position].length; x++) {
      chart_ctx.moveTo(previous_points[0], previous_points[1]);
      chart_ctx.lineTo(data_array[data_array_position][x][0], data_array[data_array_position][x][1]);
      previous_points = data_array[data_array_position][x];
    };

    chart_ctx.stroke();

  };
}

// sgr_playback_line_ctx
function draw_graph_playback_line(packet_number, chart_ctx, chart_canvas) {
  clear_chart(chart_ctx, chart_canvas);

  chart_ctx.strokeStyle = sgr_playback_line_color;

  chart_ctx.beginPath();
  chart_ctx.moveTo(packet_number, 0);
  chart_ctx.lineTo(packet_number, chart_canvas.height);

  chart_ctx.stroke();
}


function playback(frame) {
  frame_number++;

  // motionData_array
  // There is a packet for every frame
  if (motionData_array[packet_number_motionData_array].Frame_identifier == frame_number) {
    fl_suspension_position_data.innerHTML = motionData_array[packet_number_motionData_array].Suspension_position_fl;
    bl_suspension_position_data.innerHTML = motionData_array[packet_number_motionData_array].Suspension_position_rl;
    fr_suspension_position_data.innerHTML = motionData_array[packet_number_motionData_array].Suspension_position_fr;
    br_suspension_position_data.innerHTML = motionData_array[packet_number_motionData_array].Suspension_position_rr;

    packet_number_motionData_array++;
  }

  // lapData_array
  if (lapData_array[0][packet_number_lapData_array].Frame_identifier == frame_number) {
    total_time.innerHTML = intTime_to_timeTime(lapData_array[0][packet_number_lapData_array].M_currentLapTime);
    current_lap.innerHTML = intTime_to_timeTime(lapData_array[0][packet_number_lapData_array].M_currentLapTime);
    last_lap.innerHTML = intTime_to_timeTime(lapData_array[0][packet_number_lapData_array].M_lastLapTime);
    best_lap.innerHTML = intTime_to_timeTime(lapData_array[0][packet_number_lapData_array].M_bestLapTime);
    penalties.innerHTML = lapData_array[0][packet_number_lapData_array].M_penalties;
    current_sector.innerHTML = lapData_array[0][packet_number_lapData_array].M_sector;
    sector_1.innerHTML = intTime_to_timeTime(lapData_array[0][packet_number_lapData_array].M_sector1Time);
    sector_2.innerHTML = intTime_to_timeTime(lapData_array[0][packet_number_lapData_array].M_sector2Time);
    sector_3.innerHTML = intTime_to_timeTime(0);

    packet_number_lapData_array++;
  }

  // telemetryData_array
  if (telemetryData_array[packet_number_telemetryData_array].Frame_identifier == frame_number) {
    fl_tyre_pressure_data.innerHTML = telemetryData_array[packet_number_telemetryData_array].M_tyresPressure_fl;
    fl_tyre_temp_data.innerHTML = telemetryData_array[packet_number_telemetryData_array].M_tyresSurfaceTemperature_fl;
    fl_break_temp_data.innerHTML = telemetryData_array[packet_number_telemetryData_array].M_brakesTemperature_fl;
    bl_break_temp_data.innerHTML = telemetryData_array[packet_number_telemetryData_array].M_brakesTemperature_rl;
    bl_tyre_pressure_data.innerHTML = telemetryData_array[packet_number_telemetryData_array].M_tyresPressure_rl;
    bl_tyre_temp_data.innerHTML = telemetryData_array[packet_number_telemetryData_array].M_tyresSurfaceTemperature_rl;
    fr_tyre_temp_data.innerHTML = telemetryData_array[packet_number_telemetryData_array].M_tyresSurfaceTemperature_fr;
    fr_tyre_pressure_data.innerHTML = telemetryData_array[packet_number_telemetryData_array].M_tyresPressure_fr;
    fr_break_temp_data.innerHTML = telemetryData_array[packet_number_telemetryData_array].M_brakesTemperature_fr;
    br_break_temp_data.innerHTML = telemetryData_array[packet_number_telemetryData_array].M_brakesTemperature_rr;
    br_tyre_temp_data.innerHTML = telemetryData_array[packet_number_telemetryData_array].M_tyresSurfaceTemperature_rr;
    br_tyre_pressure_data.innerHTML = telemetryData_array[packet_number_telemetryData_array].M_tyresPressure_rr;


    draw_graph_playback_line(packet_number_telemetryData_array, sgr_playback_line_ctx, sgr_canvas)

    packet_number_telemetryData_array++;
  }

  // statusData_array
  if (statusData_array[packet_number_statusData_array].Frame_identifier == frame_number) {
    fl_tyre_wear_data.innerHTML = statusData_array[packet_number_statusData_array].M_tyresWear_fl;
    fl_tyre_damage_data.innerHTML = statusData_array[packet_number_statusData_array].M_tyresDamage_fl;
    bl_tyre_damage_data.innerHTML = statusData_array[packet_number_statusData_array].M_tyresDamage_rl;
    bl_tyre_wear_data.innerHTML = statusData_array[packet_number_statusData_array].M_tyresWear_rl;
    fr_tyre_wear_data.innerHTML = statusData_array[packet_number_statusData_array].M_tyresWear_fr;
    fr_tyre_damage_data.innerHTML = statusData_array[packet_number_statusData_array].M_tyresDamage_fr;
    br_tyre_damage_data.innerHTML = statusData_array[packet_number_statusData_array].M_tyresDamage_rr;
    br_tyre_wear_data.innerHTML = statusData_array[packet_number_statusData_array].M_tyresWear_rr;

    packet_number_statusData_array++;
  }

  // standings_standingsData_array
  // standings_lapDataTimes_array
  if (standings_standingsData_array[packet_number_standings_standingsData_array].Frame_identifier == frame_number) {
    for (racer_standing = 0; racer_standing < participantData_array.length; racer_standing++) {
      standing_color_list[racer_standing].innerHTML = participantData_array[standings_standingsData_array[packet_number_standings_standingsData_array].Standings[racer_standing] - 1].M_raceNumber;
      standing_name_list[racer_standing].innerHTML = participantData_array[standings_standingsData_array[packet_number_standings_standingsData_array].Standings[racer_standing] - 1].M_name;
      standing_time_list[racer_standing].innerHTML = intTime_to_timeTime(standings_lapDataTimes_array[packet_number_standings_lapDataTimes_array].Times[racer_standing - 1]);
    }

    packet_number_standings_standingsData_array++;
    packet_number_standings_lapDataTimes_array++;
  }

}

function start_playback(frame_number, play_back_speed_multiplier) {
  play_back_speed = 14.2857 * play_back_speed_multiplier;
  console.log("playback_speed", play_back_speed);
  playback_id = setInterval(playback, play_back_speed, frame_number);
}

function stop_playback() {
  clearInterval(playback_id);
  playback_id = 0;
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

      if (save_session_alert.classList.contains("hide")) {
        save_session_alert.classList.toggle('show');
        save_session_alert.classList.toggle('hide');
      }

      for (z = 0; z < data.Num_of_sessions; z++) {
        add_session_row(data.Sessions[z].Session_UID, data.Sessions[z].Session_start_time, data.Sessions[z].Session_end_time, "save_popup_body")
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


    case 40:
      console.log("packet 40 motionData recieved");
      motionData_array = data.MotionData;
      break;

    case 41:
      console.log("packet 41 sessionData recieved");
      // Get amount of laps from this to then set up everything else
      // console.log("amount of laps is:", data.SessionData[(data.SessionData).length -1].M_totalLaps)
      totalLaps = data.SessionData[(data.SessionData).length - 1].M_totalLaps;
      trackId = data.SessionData[(data.SessionData).length - 1].M_trackId;
      break;

    case 42:
      console.log("packet 42 lapData recieved");

      for (lap = 0; lap < data.LapData.length; lap++) {
        console.log("lap ", data.LapData[lap].LapNum, " is included");
        add_lap_to_lap_selector(data.LapData[lap].LapNum);
        lapData_array.push(data.LapData[lap].LapData_list);
      }

      console.log("lapData_array", lapData_array);
      break;

    case 44:
      console.log("packet 44 participantData recieved")
      participantData_dict = data.ParticipantData;
      participantData_array = data.ParticipantData;
      console.log("participantData_array", participantData_array);
      break;

    case 46:
      console.log("packet 46 telemetryData recieved");
      for (packet = 0; packet < data.TelemetryData.length; packet++) {
        speed_data_array.push([packet, data.TelemetryData[packet].M_speed * speed_multiplier]);
        gear_data_array.push([packet, data.TelemetryData[packet].M_gear * gear_multiplier]);
        rpm_data_array.push([packet, data.TelemetryData[packet].M_engineRPM * rpm_multiplier]);
      }
      draw_sgr_chart(speed_data_array, gear_data_array, rpm_data_array, sgr_ctx, sgr_canvas)

      telemetryData_array = data.TelemetryData;
      break;

    case 47:
      console.log("packet 47 statusData recieved");
      statusData_array = data.StatusData;
      console.log("statusData_array", statusData_array);
      break;

    case 48:
      console.log("packet 48 standings data recieved")
      standings_standingsData_array = data.StandingsData;
      standings_lapDataTimes_array = data.LapDataTimes;
      break;

  }
}
