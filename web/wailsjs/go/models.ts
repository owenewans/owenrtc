export namespace instances {
	
	export class OutboundSOCKS {
	    host: string;
	    port: number;
	    user?: string;
	    pass?: string;
	
	    static createFrom(source: any = {}) {
	        return new OutboundSOCKS(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.host = source["host"];
	        this.port = source["port"];
	        this.user = source["user"];
	        this.pass = source["pass"];
	    }
	}
	export class Limits {
	    traffic_limit: number;
	    speed_limit: number;
	
	    static createFrom(source: any = {}) {
	        return new Limits(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.traffic_limit = source["traffic_limit"];
	        this.speed_limit = source["speed_limit"];
	    }
	}
	export class Instance {
	    id: string;
	    name: string;
	    provider: string;
	    transport: string;
	    room_id: string;
	    limits: Limits;
	    outbound?: OutboundSOCKS;
	    created_at: number;
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new Instance(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.provider = source["provider"];
	        this.transport = source["transport"];
	        this.room_id = source["room_id"];
	        this.limits = this.convertValues(source["limits"], Limits);
	        this.outbound = this.convertValues(source["outbound"], OutboundSOCKS);
	        this.created_at = source["created_at"];
	        this.status = source["status"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	

}

export namespace mode {
	
	export class Mode {
	    Kind: string;
	    Port: number;
	    PublicIP: string;
	
	    static createFrom(source: any = {}) {
	        return new Mode(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Kind = source["Kind"];
	        this.Port = source["Port"];
	        this.PublicIP = source["PublicIP"];
	    }
	}

}

