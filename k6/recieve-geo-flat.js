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
let keySeqLen = 100


export let options = {
    tags: {
        "name": "receive geo"
    },
    setupTimeout: "30m",
//   minIterationDuration: "100ms"
};

export function setup() {
    var keys = clusterKeys()
    var keySeqs = []
    for (var iter=0;iter<__ENV["K6_ITERATIONS"];iter++){
        keySeqs[iter] = genKeySeq(keySeqLen, keys)
    }
    return keySeqs
}

function keyCompare(k0s, k1s){
    var k0 = JSON.parse(k0s)
    var k1 = JSON.parse(k1s)
    if (k0.Lon > k1.Lon){
        return true
    }
    if (k0.Lat > k1.Lat){
        return true
    }
    return false
}

function clusterKeys(){
    var cluster = http.get(p_hosts[0] + "/list");
    var body = JSON.parse(cluster.body)
    var keys = []

    Object.getOwnPropertyNames(body).forEach(function(node) {
        keys = keys.concat(body[node])
    })

    keys.sort(keyCompare)

    return keys
}

function genKeySeq(count, keys) {
    var start =  Math.floor(Math.random() * (keys.length-count))
    var key = ""

    for (var iter = start; iter < count+start; iter++) {
        key += "/"+keys[iter%keys.length]
        // console.log(keys[iter%keys.length])
    }
    return key
}

export default function(keySeqs) {
    let maxIter = __ENV["K6_ITERATIONS"] / __ENV["K6_VUS"]
    let idx = ((maxIter * (__VU - 1)) + __ITER) % (keySeqs.length)

    // console.log(key, idx, __ITER, __VU, maxIter, points.length)
    // console.log(key)
    var res = http.get(p_hosts[__VU % p_hosts.length] + "/get"+keySeqs[idx], null, {tags: {name: 'get_download_geo'}});
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
