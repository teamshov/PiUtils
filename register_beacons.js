// Packest from the Estimote family (Telemetry, Connectivity, etc.) are
// broadcast as Service Data (per "§ 1.11. The Service Data - 16 bit UUID" from
// the BLE spec), with the Service UUID 'fe9a'.
var ESTIMOTE_SERVICE_UUID = 'fe9a';

// Once you obtain the "Estimote" Service Data, here's how to check if it's
// a Telemetry packet, and if so, how to parse it.
function parseEstimoteTelemetryPacket(data) { // data is a 0-indexed byte array/buffer
  // byte 0, lower 4 bits => frame type, for Telemetry it's always 2 (i.e., 0b0010)
  var frameType = data.readUInt8(0) & 0b00001111;
  var ESTIMOTE_FRAME_TYPE_TELEMETRY = 2;
  if (frameType != ESTIMOTE_FRAME_TYPE_TELEMETRY) { return; }

  // byte 0, upper 4 bits => Telemetry protocol version ("0", "1", "2", etc.)
  var protocolVersion = (data.readUInt8(0) & 0b11110000) >> 4;
  // this parser only understands version up to 2
  // (but at the time of this commit, there's no 3 or higher anyway :wink:)
  if (protocolVersion > 2) { return; }

  // bytes 1, 2, 3, 4, 5, 6, 7, 8 => first half of the identifier of the beacon
  var shortIdentifier = data.toString('hex', 1, 9);
  return shortIdentifier;
}

// example how to scan & parse Estimote Telemetry packets with noble

var noble = require('noble');

noble.on('stateChange', function(state) {
  console.log('state has changed', state);
  if (state == 'poweredOn') {
    var serviceUUIDs = [ESTIMOTE_SERVICE_UUID]; // Estimote Service
    var allowDuplicates = true;
    noble.startScanning(serviceUUIDs, allowDuplicates, function(error) {
      if (error) {
        console.log('error starting scanning', error);
      } else {
        console.log('started scanning');
      }
    });
  }
});

noble.on('discover', function(peripheral) {
  var serviceData = peripheral.advertisement.serviceData.find(function(el) {
    return el.uuid == ESTIMOTE_SERVICE_UUID;
  });

  if (serviceData === undefined) { return; }
  var data = serviceData.data;

  var sid = parseEstimoteTelemetryPacket(data);
    if (sid) { handle(sid) }
});

var readline = require("readline")

const r1 = readline.createInterface({
    input: 
})

function handle(sid) {
    
}
