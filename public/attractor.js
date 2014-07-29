function find_angle(p0,p1,c) {
	var p0c = Math.sqrt(Math.pow(c.x-p0.x,2)+
                    Math.pow(c.y-p0.y,2)); // p0->c (b)   
	var p1c = Math.sqrt(Math.pow(c.x-p1.x,2)+
                    Math.pow(c.y-p1.y,2)); // p1->c (a)
	var p0p1 = Math.sqrt(Math.pow(p1.x-p0.x,2)+
                     Math.pow(p1.y-p0.y,2)); // p0->p1 (c)
	return Math.acos((p1c*p1c+p0c*p0c-p0p1*p0p1)/(2*p1c*p0c));
}

function HSVtoRGB(h, s, v) {
	var r, g, b, i, f, p, q, t;
	if (h && s === undefined && v === undefined) {
    	s = h.s, v = h.v, h = h.h;
	}
    i = Math.floor(h * 6);
	f = h * 6 - i;
    p = v * (1 - s);
	q = v * (1 - f * s);
    t = v * (1 - (1 - f) * s);
	switch (i % 6) {
    	case 0: r = v, g = t, b = p; break;
        case 1: r = q, g = v, b = p; break;
	    case 2: r = p, g = v, b = t; break;
    	case 3: r = p, g = q, b = v; break;
        case 4: r = t, g = p, b = v; break;
	    case 5: r = v, g = p, b = q; break;
	}
	return {
    	r: Math.floor(r * 255),
    	g: Math.floor(g * 255),
    	b: Math.floor(b * 255)
	};
}

var params =[
	[-2.951292,1.18775,0.517396,1.090625,300000],
	[-0.9669180,2.879879,0.76145,0.744728,200000], // King
	[-2.905148,-2.030427,1.44055,0.70307,300000], // lords
	[2.668752,1.225105,0.709998,0.637272,200000], // lords3
	[2.73394,1.369945,1.471923,0.869182, 300000], //lords5
	[1.008118,2.65392, 0.599124,0.6507, 200000], // lords6
	[-2.767266,-0.633839,1.352107,0.705481, 300000], // serfs1
	[-2.164647,-0.641713,1.277032,1.003342,50000, ],// serfs3
];

function drawImage(canvas, ctx){
	var countPoints = 0;
	var selectedParam = Math.floor(Math.random() * params.length)

	var a = params[selectedParam][0];
	var b = params[selectedParam][1];
	var c = params[selectedParam][2];
	var d = params[selectedParam][3];
	var maxPoints = params[selectedParam][4];

	var lastx = 1;
	var lasty = 1;
	var prevx = 1;
	var prevy = 1;

	var interval = setInterval(function(){
    	var h = canvas.width ;
    	var w = canvas.height;
  
    	for( var i =0; i< 2000; ++i){
    		countPoints++
			var x = Math.sin(lasty*b) + c*Math.sin(lastx*b);
			var y = Math.sin(lastx*a) + d*Math.sin(lasty*a);

			var angle = find_angle({x:prevx, y:prevy}, {x:x,y:y}, {x:lastx, y:lasty})
			
			var rgb = HSVtoRGB(angle / Math.PI *0.3 + 0.65, 0.6, 0.3);
			ctx.fillStyle = "rgba("+rgb.r+"," + rgb.g + ", " + rgb.b + ",0.2)";

			lastx = x;
			lasty = y;

    		var fx = (x + 6)/6.0 * w - w/2;
			var fy = (y+6)/6.0 * h - h/2;
        	ctx.fillRect (fx,fy, 1, 1);
    	}
    	if(countPoints > maxPoints){
    		clearInterval(interval);
    	}
	},100);

	return interval
}