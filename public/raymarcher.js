function length(p){
	var r = Math.sqrt(p[0]*p[0] + p[1]*p[1] + p[2]*p[2]);
	return r;
}
function abs(p){
	return [Math.abs(p[0]), Math.abs(p[1]), Math.abs(p[2])];
}
function sum(p1,p2){
	return [p1[0]+p2[0], p1[1]+p2[1], p1[2]+p2[2]];
}
function minus(p1,p2){
	return [p1[0]-p2[0], p1[1]-p2[1], p1[2]-p2[2]];	
}
function max(p1,p2){
	return [Math.max(p1[0],p2[0]),Math.max(p1[1],p2[1]),Math.max(p1[2],p2[2])];
}
function dot(p1,p2){
	return p1[0]*p2[0]+p1[1]*p2[1]+p1[2]*p2[2];
}
function mul(p1,p2){
	return [p1[0]*p2[0], p1[1]*p2[1], p1[2]*p2[2]];
}
function dist(p1,p2){
	return length(minus(p1,p2))
}
function scale(p,s){
	return [p[0]*s, p[1]*s,p[2]*s]
}
function normalized(p) {
	var l = length(p);
	if(l < 0.001){
		l = 0.001;
	}
	return scale(p, 1/l);
}

function sphereDist(center, radius, p) {
	return dist(center,p) - radius;
}

function boxDist(center, sides,p )
{
  var d = minus(abs(minus(p,center)), sides);
  return Math.min(Math.max(d[0],Math.max(d[1],d[2])),0.0) + length(max(d,[0,0,0]));
}

// From http://iquilezles.org/www/articles/distfunctions/distfunctions.htm
function opUnion(/*args*/){
    return Math.min.apply(null,arguments);
}
function opS(d1, d2){
    return Math.max(-d1,d2);
}
function opIntersect(/*args*/){
    return Math.max.apply(null,arguments);
}

function gradient(func, p){
	var eps = 0.01;
	return normalized([
		func(sum(p,[eps,0,0])) - func(minus(p,[eps,0,0])),
		func(sum(p,[0,eps,0])) - func(minus(p,[0,eps,0])),
		func(sum(p,[0,0,eps])) - func(minus(p,[0,0,eps]))]);
}

function ao(func, p){
	var normal = gradient(func,p);
	var aosum = 0;
	var n = 10;
	var k = 2.5 / n;
	var delta = 0.05;
	for(var i=0;i<n;++i){
		aosum += i * delta - func( sum(p , scale(normal, i*delta) )) / Math.pow(2,i);
	}
	
	return 1 - k * aosum;
}

function softshadow(func, ro, rd, mint, k){
    var res = 1.0;
    var t = mint;
	var h = 1.0;
    for( var i=0; i<50; i++ )
    {
        var h = func(sum(ro, scale(rd,t)));
        res = Math.min( res, k*h/t );
		t += Math.min(Math.max(h, 0.02), 2.0/50.0);
		if( res<0.01 ) break;
    }
    return Math.min(Math.max(res,0.0),1.0);
}

function dist1(p){
	var d1 = sphereDist([3,-2.5,1], 1.5, p)+ 0.3*Math.sin(5*p[2]+3.2*p[0])*Math.cos(5*p[2]*p[1]*p[0]);
	var d2 = sphereDist([-3,-2,1], 1.5, p) + 0.2 * Math.sin(3*p[0] + 4*p[1]) * Math.cos(4*p[2]*p[0]) * Math.sin(5*p[2]*p[0]);
	
	var d3 = boxDist([4,4,10], [4,3,6],p) + 0.2*Math.sin(13.2*Math.cos(p[0]*p[1]- p[0]*p[2]));
	return opUnion(d1,d2,d3);
}


function drawImage(canvas, ctx){
	var sceneDist = dist1;
	var countPoints = 0;
	var maxPoints = 100000;

	var lightDir = normalized([-0.5,0.3,-0.5]);
	var interval = setInterval(function(){
    	var w = canvas.width ;
    	var h = canvas.height;
  
    	for( var i =0; i< 4000; ++i){
    		
			var x = (2 * Math.random()) - 1.0;
			var y = (2 * Math.random()) - 1.0;

			var p = [0,0,-5];

			var d = normalized([
				1.4 * x,
				1.4 * y,
				1]);
			
			var r = Math.min(w,h) * 0.8;
			x = x * r + w/2;
			y = -y * r + h/2;
			
			var drawn = false
			for( var steps =0; steps < 150; ++steps){

				var distance = sceneDist(p);
				if(distance < 0.01){
					var normal = gradient(sceneDist,p);
					var occ = ao(sceneDist,p);
					
					var lambert = occ + 0.6*Math.max(0,dot(normal, lightDir));
					lambert = lambert * (0.6 + 0.4 * softshadow(sceneDist, sum(p,scale(normal,0.04)), lightDir, 0.005, 64.0 ));

					var darkNormal = mul(normal, [0.1,0,0.1]);
					
					var illum = sum(darkNormal,[lambert,lambert,lambert]);

					var r = Math.floor(illum[0]*255);
					var g = Math.floor(illum[1]*255);
					var b = Math.floor(illum[2]*255);

					ctx.fillStyle = "rgb(" + r + "," + g+","+b+")";
        			ctx.fillRect (x,y, 1, 1);
        			drawn = true;
        			countPoints++;
        			break;
				}
				if(distance > 30){
					break;
				}
				p = sum(p, scale(d, distance * 0.98));
			}
				

    	}
    	if(countPoints > maxPoints){
    		clearInterval(interval);
    	}
	},100);

	return interval
}