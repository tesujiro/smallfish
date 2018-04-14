var map;

var changeMap = function(lat,lng){
	map = new google.maps.Map(document.getElementById('map'), {
		center: {lat: lat, lng: lng},
		zoom: 14
	});
}

var doPost = function(jsonArray){
	console.log('doPost:'+jsonArray);
	if (jsonArray.length==0) return;
	//console.log("jsonArray.length=="+jsonArray.length);
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
	post          : function() {
		doPost(this.json);
		this.clearJson();
		this.postTimer=setTimeout(this.post.bind(this), 5000);
	}
}

var changeLocation = function() {
	/*
	var json=[];
	var clearJson = function() {
		json=[];
	}
	var pushJson = function(id,time,lat,lng){
		json.push({
			"consumerId"	: id ,
			"timestamp"	: time ,
			"latitude"	: lat ,
			"longtitude"	: lng
		});
	}
	var postTimer=0;
	var stopPostTimer = function() {
		clearTimeout(postTimer);
		postTimer=0;
	}
	var postGeoInfo = function() {
		doPost(json);
		clearJson();
		postTimer=setTimeout(function(){postGeoInfo()}, 5000);
	}
	*/

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
			//pushJson(1 ,new Date() , event.latLng.lat() , event.latLng.lng());
			info.pushJson(1 ,new Date() , event.latLng.lat() , event.latLng.lng());
			//console.log("info:"+info.json.length);
		});
		//postGeoInfo();
		console.log("before post info:"+info.json.length);
		info.post();
	};
	//stopPostTimer();
	info.stopPostTimer();
	navigator.geolocation.getCurrentPosition(geoSuccess);
};
