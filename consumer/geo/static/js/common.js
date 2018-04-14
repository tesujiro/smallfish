var map;

var changeMap = function(lat,lng){
	map = new google.maps.Map(document.getElementById('map'), {
		center: {lat: lat, lng: lng},
		zoom: 14
	});
}

var doPost = function(jsonArray){
	console.log('doPost:'+jsonArray.length);
	if (jsonArray.length==0) return;
	var req = new XMLHttpRequest();
	req.onreadystatechange = function() {
		if (req.readyState == 4) { // finished sending
			console.log("req.status="+req.status);
			if (req.status == 200) {
				console.log(req.responseText);
			}
		}else{
			console.log("通信中...");
		}
	}
	req.open('POST', '/consumer/GeoCollection', true);
	req.setRequestHeader("Content-type", "application/json");
	var parameters = JSON.stringify(jsonArray);
	req.send(parameters);
}

function geoInfo() {
	this.json = [];
	this.postTimer = 0;
	//this.Timeout = 5000; // 5 seconds
	this.Timeout = 60000; // 60 seconds
};
geoInfo.prototype = {
	//json      : [] ,
	clearJson : function() {
		this.json=[];
	},
	pushJson  : function(id,time,lat,lng){
		this.json.push({
			"consumerId"	: id ,
			"timestamp"	: time ,
			"latitude"	: lat ,
			"longtitude"	: lng
		});
		//console.log("pushJson:"+this.json.length+" :"+this.json[this.json.length-1]);
	},
	//postTimer     : 0,
	stopPostTimer : function() {
		clearTimeout(this.postTimer);
		this.postTimer=0;
	},
	startPost: function(){
		this.stopPostTimer(); // avoid duplicate timer
		this.postTimer=setTimeout(this.post.bind(this), this.Timeout);
	},
	post          : function() {
		doPost(this.json);
		this.clearJson();
		this.postTimer=setTimeout(this.post.bind(this), this.Timeout);
	}
}

var changeLocation = function() {

	var info = new geoInfo();
	var currentPos;
	var geoSuccess = function(position) {
		currentPos = position;
		console.log('Lat=' + currentPos.coords.latitude + ' Lng=' + currentPos.coords.longitude);
		document.getElementById('currentLat').innerHTML = currentPos.coords.latitude;
		document.getElementById('currentLon').innerHTML = currentPos.coords.longitude;
		changeMap(currentPos.coords.latitude,currentPos.coords.longitude);
		// Update lat/long value of div when anywhere in the map is clicked
		google.maps.event.addListener(map,'click',function(event) {
			document.getElementById('currentLat').innerHTML = event.latLng.lat();
			document.getElementById('currentLon').innerHTML = event.latLng.lng();
			info.pushJson(1 ,new Date() , event.latLng.lat() , event.latLng.lng());
		});
		info.startPost();
	};
	navigator.geolocation.getCurrentPosition(geoSuccess);
};

function autoPost() {
	this.info = new geoInfo();
	//this.Timeout = 1000; // 1 second
	this.Timeout = 10000; // 10 seconds
	this.postTimer = 0;
};

autoPost.prototype = {
	stopAutoPostTimer : function() {
		clearTimeout(this.postTimer);
		this.postTimer = 0;
	},
	autoGeo : function() {
		var geoSuccess = function(position) {
			currentPos = position;
			console.log('Lat=' + currentPos.coords.latitude + ' Lng=' + currentPos.coords.longitude);
			document.getElementById('currentLat').innerHTML = currentPos.coords.latitude;
			document.getElementById('currentLon').innerHTML = currentPos.coords.longitude;
			this.info.pushJson(1 ,new Date() , currentPos.coords.latitude ,currentPos.coords.longitude); 
			this.postTimer=setTimeout(this.autoGeo.bind(this), this.Timeout);
		};
		navigator.geolocation.getCurrentPosition(geoSuccess.bind(this));
	},
	start : function() {
		this.stopAutoPostTimer(); // avoid duplicate timer
		this.info.stopPostTimer();
		this.postTimer = setTimeout(this.autoGeo.bind(this), this.Timeout);
		this.info.startPost();
	},
	stop : function() {
		this.stopAutoPostTimer();
		this.info.stopPostTimer();
	}
}

var post ;
var startAutoPost = function() {
	post = new autoPost();
	post.start()
};
var stopAutoPost = function() {
	post.stop()
};
