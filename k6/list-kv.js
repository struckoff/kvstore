import http from "k6/http";
import { sleep } from "k6";
import { group } from "k6";
import { check } from "k6";


// let p_host = "http://ip-172-31-42-150.eu-central-1.compute.internal:47375"
let p_host = "http://localhost:9190"
// let dataFile = open("/home/struckoff/Documents/alldata/cut_10k.csv");

export let options = {
    tags: {
        "name": "list k/v"
    },
//   minIterationDuration: "100ms"
};

export default function() {
    var res = http.get(p_host + "/list", {tags: {name: 'get_list_kv'}});
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