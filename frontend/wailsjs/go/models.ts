export namespace config {
	
	export class Config {
	    appName: string;
	    version: string;
	    ollamaUrl: string;
	    ollamaModel: string;
	    kakaoEnabled: boolean;
	    kakaoPort: number;
	    kakaoWebhookPath: string;
	    kakaoDmPolicy: string;
	    kakaoAllowFrom: string[];
	    kakaoSystemPrompt: string;
	    kakaoModel: string;
	    relayUrl: string;
	    relayToken: string;
	    authProvider: string;
	    accessToken: string;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.appName = source["appName"];
	        this.version = source["version"];
	        this.ollamaUrl = source["ollamaUrl"];
	        this.ollamaModel = source["ollamaModel"];
	        this.kakaoEnabled = source["kakaoEnabled"];
	        this.kakaoPort = source["kakaoPort"];
	        this.kakaoWebhookPath = source["kakaoWebhookPath"];
	        this.kakaoDmPolicy = source["kakaoDmPolicy"];
	        this.kakaoAllowFrom = source["kakaoAllowFrom"];
	        this.kakaoSystemPrompt = source["kakaoSystemPrompt"];
	        this.kakaoModel = source["kakaoModel"];
	        this.relayUrl = source["relayUrl"];
	        this.relayToken = source["relayToken"];
	        this.authProvider = source["authProvider"];
	        this.accessToken = source["accessToken"];
	    }
	}

}

export namespace main {
	
	export class ChatMessage {
	    role: string;
	    content: string;
	    timestamp: number;
	
	    static createFrom(source: any = {}) {
	        return new ChatMessage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.role = source["role"];
	        this.content = source["content"];
	        this.timestamp = source["timestamp"];
	    }
	}
	export class HWPParseResult {
	    success: boolean;
	    text?: string;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new HWPParseResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.text = source["text"];
	        this.error = source["error"];
	    }
	}
	export class KakaoConfigJS {
	    enabled: boolean;
	    port: number;
	    webhookPath: string;
	    dmPolicy: string;
	    allowFrom: string[];
	    systemPrompt: string;
	    model: string;
	
	    static createFrom(source: any = {}) {
	        return new KakaoConfigJS(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enabled = source["enabled"];
	        this.port = source["port"];
	        this.webhookPath = source["webhookPath"];
	        this.dmPolicy = source["dmPolicy"];
	        this.allowFrom = source["allowFrom"];
	        this.systemPrompt = source["systemPrompt"];
	        this.model = source["model"];
	    }
	}
	export class KakaoStatus {
	    running: boolean;
	    enabled: boolean;
	    port: number;
	    webhookPath: string;
	
	    static createFrom(source: any = {}) {
	        return new KakaoStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.running = source["running"];
	        this.enabled = source["enabled"];
	        this.port = source["port"];
	        this.webhookPath = source["webhookPath"];
	    }
	}
	export class SystemInfo {
	    os: string;
	    arch: string;
	    memoryGb: number;
	
	    static createFrom(source: any = {}) {
	        return new SystemInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.os = source["os"];
	        this.arch = source["arch"];
	        this.memoryGb = source["memoryGb"];
	    }
	}

}

