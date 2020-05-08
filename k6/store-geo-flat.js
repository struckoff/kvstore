import http from "k6/http";
import { sleep } from "k6";
import { group } from "k6";
import { check } from "k6";


// let p_host = "http://ip-172-31-42-150.eu-central-1.compute.internal:47375"
let p_hosts = [
    "http://localhost:9190",
    "http://localhost:9191",
    "http://localhost:9192"
]
// let dataFile = open("/home/struckoff/Documents/alldata/cut_10k.csv");

let minLat = -90
let maxLat = 90
let minLon = -180
let maxLon = 180

let step = .5

export let options = {
    tags: {
        "name": "store geo"
    },
//   minIterationDuration: "100ms"
};

function randomPoint(){
    // var timeInMs = Date.now() - Math.floor(Math.random());
    var timeInMs = Date.now();
    var val = Math.random() + 10
    return timeInMs + ';' + val
}

function genData(){
    let iter = 0;
    let data = [];
    for (var lat = minLat; lat <maxLat; lat+= step) {
        for (var lon = minLon; lon <maxLon; lon+= step){
            let d = {"Lon":lon, "Lat":lat}
            data[iter] = JSON.stringify(d)
            iter++
        }
    }
    return data
}

export function setup() {
    var points = genData()
    // var index = 0
    // while (index < 10000){
    //     yield points[index]
    //     index++
    // }
    return points
}

export default function(points) {
    let maxIter = __ENV["K6_ITERATIONS"] / __ENV["K6_VUS"]
    let idx = ((maxIter * (__VU - 1)) + __ITER) % (points.length)
    let key = points[idx]

    // console.log(key, idx, __ITER, __VU, maxIter, points.length)
    var res = http.post(p_hosts[__VU % p_hosts.length] + "/put/"+key, randomPoint(), {tags: {name: 'post_upload_geo'}});
    if (res.status >= 400){
        console.error(res.body)
    }
    check(res, {
        "is status OK": (r) => r.status < 400,
        "is status not 404": (r) => r.status != 404,
        "is status not 403": (r) => r.status != 403,
        "is status not 500": (r) => r.status != 500,
        "is status not 503": (r) => r.status != 503,
    });
}