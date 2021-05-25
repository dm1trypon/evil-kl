package schemas

// SchemasMap - API method schemas
var SchemasMap = map[string]string{
	"getKeyloggerData": `
{
	"type": "object",
	"required": [
		"method"
	],
	"properties": {
		"method": {
			"type": "string",
			"minLength": 1,
			"enum": [
				"getKeyloggerData"
			]
		}
	},
	"additionalProperties": true
}
`,
	"getLogs": `
{
	"type": "object",
	"required": [
		"method"
	],
	"properties": {
		"method": {
			"type": "string",
			"minLength": 1,
			"enum": [
				"getLogs"
			]
		}
	},
	"additionalProperties": true
}`,
}
