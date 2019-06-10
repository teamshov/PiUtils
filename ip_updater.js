const os = require('os');
const localip = require('local-ip');
const iface = 'wlan0';
const nano = require('nano')('http://admin:seniorshov@omaraa.ddns.net:5984');
const pies = nano.db.use('pies');

localip(iface, function(err, res) {
	if(err){return;}
	
    pies.get(os.hostname(), function(err, pi) {
	pies.insert({ip:res, _rev:pi._rev}, os.hostname());
	console.log(os.hostname()+" has been updated with ip: "+res);
	});
});



