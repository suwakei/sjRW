package sjrw

import (
	"log"
	"testing"
)


func BenchmarkReadAsStr(b *testing.B) {
	var jsonPath string = "./testdata/readtest.json"
	var sj SjReader
	for i := 0; i < 100; i++ {
	_, err := sj.ReadAsStrFrom(jsonPath)
	if err != nil {
		log.Fatal(err)
	}
}
}


func BenchmarkReadAsBytes(b *testing.B) {
	var jsonPath string = "./testdata/readtest.json"
	var sj SjReader
	for i := 0; i < 100; i++ {
	_, err := sj.ReadAsBytesFrom(jsonPath)
	if err != nil {
		log.Fatal(err)
	}
}
}


func TestReadAsStr(t *testing.T) {
	var jsonPath1 string = "./testdata/readtest.json"
	var sj1 SjReader
	input1, _ := sj1.ReadAsStrFrom(jsonPath1)

	var jsonPath2 string = "./testdata/readtest2.json"
	var sj2 SjReader
	input2, _ := sj2.ReadAsStrFrom(jsonPath2)
	tests := []struct{
		input string
		expected string
	}{
		{
			input1,
			expected1,
		},
		{
			input2,
			expected2,
		},
	}

	for _, tt := range tests {
		testReadAsStrFrom(t, tt.input, tt.expected)
	}
}

func testReadAsStrFrom(t *testing.T, ReadAsStrResult, ReadAsStrExpected string) bool {
	if ReadAsStrResult != ReadAsStrExpected {
		t.Errorf("these values are not same")
		return false
	}
	return true
}





const expected1 string = `
[
    {
      "_id": "672d31b26f1316908fa81a41",
      "index": 0,
      "guid": "126bf441-05a3-4b3e-9868-43827b2054c4",
      "isActive": false,
      "balance": "$1,509.41",
      "picture": "http://placehold.it/32x32",
      "age": 26,
      "eyeColor": "green",
      "name": "Sue Irwin",
      "gender": "female",
      "company": "VURBO",
      "email": "sueirwin@vurbo.com",
      "phone": "+1 (829) 544-2803",
      "address": "196 Pierrepont Street, Bartonsville, Idaho, 2203",
      "about": "Duis nisi Lorem occaecat do eu fugiat consectetur. Reprehenderit ut magna velit est reprehenderit Lorem. Excepteur consequat velit enim veniam quis velit velit enim aliquip nisi commodo ex. Pariatur consequat laboris amet fugiat nulla quis duis irure proident duis. Elit irure officia consequat reprehenderit commodo ad. Reprehenderit amet pariatur voluptate laboris dolor et veniam non ex. Consectetur ipsum in sunt irure cupidatat voluptate id ipsum.\r\n",
      "registered": "2024-01-07T08:57:16 -09:00",
      "latitude": 48.927862,
      "longitude": -79.629704,
      "tags": [
        "consequat",
        "ex",
        "tempor",
        "dolor",
        "nisi",
        "occaecat",
        "quis"
      ],
      "friends": [
        {
          "id": 0,
          "name": "Michael Rojas"
        },
        {
          "id": 1,
          "name": "Rodgers Pennington"
        },
        {
          "id": 2,
          "name": "Cardenas Monroe"
        }
      ],
      "greeting": "Hello, Sue Irwin! You have 7 unread messages.",
      "favoriteFruit": "banana"
    },
    {
      "_id": "672d31b2240100375952cb1e",
      "index": 1,
      "guid": "98329871-09ba-4d38-99e3-8d808395bdfa",
      "isActive": false,
      "balance": "$2,832.86",
      "picture": "http://placehold.it/32x32",
      "age": 27,
      "eyeColor": "blue",
      "name": "Noemi Hays",
      "gender": "female",
      "company": "OBLIQ",
      "email": "noemihays@obliq.com",
      "phone": "+1 (940) 540-2100",
      "address": "368 Suydam Place, Ezel, Ohio, 7763",
      "about": "Pariatur dolore commodo ex aliqua tempor qui sit. Dolor incididunt nulla anim occaecat excepteur consectetur commodo officia voluptate tempor voluptate eiusmod officia ut. Officia labore aliquip fugiat amet aliquip excepteur et qui et laboris aliquip sunt occaecat. Dolore sint amet Lorem ea. Ipsum magna esse amet culpa sunt est nostrud id ut nostrud pariatur eu sit anim. Occaecat Lorem aute elit reprehenderit est reprehenderit velit exercitation qui amet. Enim consequat incididunt est velit ad.\r\n",
      "registered": "2023-03-12T06:12:32 -09:00",
      "latitude": 44.210044,
      "longitude": -165.285857,
      "tags": [
        "exercitation",
        "pariatur",
        "amet",
        "qui",
        "deserunt"
      ],
      "friends": [
        {
          "id": 0,
          "name": "Christa Cole"
        },
        {
          "id": 1,
          "name": "Nona Knowles",
          "age": 32,
          "sex": "male",
          "status": "student"

        },
        {
          "id": 2,
          "name": "Keisha Mosley"
        }
      ],
      "greeting": "Hello, Noemi Hays! You have 9 unread messages.",
      "favoriteFruit": "strawberry"
    },
    {
      "_id": "672d31b23cbd5ea5d1a839c6",
      "index": 2,
      "guid": "f6cf6351-3bdd-4a3f-a353-f612e2befa91",
      "isActive": true,
      "balance": "$3,136.55",
      "picture": "http://placehold.it/32x32",
      "age": 23,
      "eyeColor": "brown",
      "name": "Parks Buck",
      "gender": "male",
      "company": "ZOSIS",
      "email": "parksbuck@zosis.com",
      "phone": "+1 (842) 554-2877",
      "address": "283 Kenmore Terrace, Shawmut, District Of Columbia, 1837",
      "about": "Proident sint id fugiat minim ut mollit nostrud magna eiusmod commodo. Laborum culpa amet ea occaecat. Do tempor labore irure officia dolor magna ullamco reprehenderit tempor. Deserunt qui dolore reprehenderit amet dolor adipisicing fugiat exercitation ad voluptate aute. Culpa amet culpa ea dolor cillum exercitation pariatur sit eu cillum dolor fugiat.\r\n",
      "registered": "2021-11-06T03:15:03 -09:00",
      "latitude": -57.562182,
      "longitude": 29.246979,
      "tags": [
        "cupidatat",
        "do",
        "Lorem",
        "cupidatat",
        "quis",
        "reprehenderit",
        "do"
      ],
      "friends": [
        {
          "id": 0,
          "name": "Alma Neal"
        },
        {
          "id": 1,
          "name": "Trudy Crawford"
        },
        {
          "id": 2,
          "name": "Bernadette Whitaker"
        }
      ],
      "greeting": "Hello, Parks Buck! You have 10 unread messages.",
      "favoriteFruit": "apple"
    },
    {
      "_id": "672d31b2e474ce8b790e060e",
      "index": 3,
      "guid": "879c9095-f972-4519-af5e-e84878cc8ae9",
      "isActive": false,
      "balance": "$2,673.28",
      "picture": "http://placehold.it/32x32",
      "age": 29,
      "eyeColor": "blue",
      "name": "Hester Petersen",
      "gender": "male",
      "company": "ZERBINA",
      "email": "hesterpetersen@zerbina.com",
      "phone": "+1 (962) 504-3557",
      "address": "181 Beard Street, Wedgewood, Northern Mariana Islands, 1108",
      "about": "Ut et quis enim dolore mollit. Reprehenderit excepteur esse Lorem ea veniam voluptate non tempor. Deserunt qui cupidatat est sunt. Sit aliqua nulla mollit dolore cillum commodo magna. Dolor non esse aute duis commodo Lorem dolor quis sit.\r\n",
      "registered": "2016-06-21T10:37:26 -09:00",
      "latitude": 81.428735,
      "longitude": 62.308805,
      "tags": [
        "voluptate",
        "quis",
        "ea",
        "veniam",
        "sunt",
        "deserunt",
        "commodo"
      ],
      "friends": [
        {
          "id": 0,
          "name": "Lori Gutierrez"
        },
        {
          "id": 1,
          "name": "Nola Griffin"
        },
        {
          "id": 2,
          "name": "Corrine Gross"
        }
      ],
      "greeting": "Hello, Hester Petersen! You have 10 unread messages.",
      "favoriteFruit": "strawberry"
    },
    {
      "_id": "672d31b2aacc9b71b972f01f",
      "index": 4,
      "guid": "33fadf23-b3d6-4710-8be5-07d6804ca411",
      "isActive": false,
      "balance": "$2,411.09",
      "picture": "http://placehold.it/32x32",
      "age": 32,
      "eyeColor": "brown",
      "name": "Margret Lloyd",
      "gender": "female",
      "company": "ORBAXTER",
      "email": "margretlloyd@orbaxter.com",
      "phone": "+1 (884) 556-3656",
      "address": "977 Schaefer Street, Olney, Connecticut, 5275",
      "about": "Minim irure nulla duis mollit exercitation nisi aliqua tempor consectetur magna. Occaecat velit do veniam ipsum magna veniam elit. Voluptate nisi Lorem ullamco et. Voluptate elit consectetur enim Lorem in ipsum quis. In Lorem exercitation ullamco nulla magna laboris fugiat labore amet. Do mollit labore cillum est veniam amet mollit reprehenderit. Laborum enim ullamco ad ea consectetur ex occaecat ut dolore.\r\n",
      "registered": "2021-09-28T05:31:32 -09:00",
      "latitude": -80.987436,
      "longitude": -126.500677,
      "tags": [
        "ea",
        "quis",
        "aliquip",
        "aliquip",
        "sint",
        "enim",
        "quis"
      ],
      "friends": [
        {
          "id": 0,
          "name": "Daniel Pace"
        },
        {
          "id": 1,
          "name": "Simmons Nunez"
        },
        {
          "id": 2,
          "name": "Roach Manning"
        }
      ],
      "greeting": "Hello, Margret Lloyd! You have 2 unread messages.",
      "favoriteFruit": "strawberry"
    },
    {
      "_id": "672d31b224f138c8a64311fc",
      "index": 5,
      "guid": "2da6913b-6155-4967-988f-294cc54f64de",
      "isActive": true,
      "balance": "$1,152.90",
      "picture": "http://placehold.it/32x32",
      "age": 35,
      "eyeColor": "green",
      "name": "Velasquez Greer",
      "gender": "male",
      "company": "ACRODANCE",
      "email": "velasquezgreer@acrodance.com",
      "phone": "+1 (826) 408-2675",
      "address": "588 Malbone Street, Catherine, Maine, 8474",
      "about": "Quis incididunt consectetur pariatur ipsum deserunt ea nisi ullamco. Minim in occaecat adipisicing nisi id labore Lorem cillum cupidatat fugiat dolor commodo. Sunt pariatur ipsum dolor elit aliquip laborum veniam aliqua consectetur sit sint nisi labore.\r\n",
      "registered": "2021-10-10T10:32:14 -09:00",
      "latitude": -65.290247,
      "longitude": 116.676104,
      "tags": [
        "irure",
        "laboris",
        "consectetur",
        "deserunt",
        "aliqua",
        "deserunt",
        "fugiat"
      ],
      "friends": [
        {
          "id": 0,
          "name": "Berry Zamora"
        },
        {
          "id": 1,
          "name": "Viola Myers"
        },
        {
          "id": 2,
          "name": "Leanna Pratt"
        }
      ],
      "greeting": "Hello, Velasquez Greer! You have 10 unread messages.",
      "favoriteFruit": "banana"
    }
  ]`



const expected2 string = `
[
	{
	  "_id": "672d31b26f1316908fa81a41",
	  "index": 0,
	  "guid": "126bf441-05a3-4b3e-9868-43827b2054c4",
	  "isActive": false,
	  "balance": "$1,509.41",
	  "picture": "http://placehold.it/32x32",
	  "age": 26,
	  "eyeColor": "green",
	  "name": "Sue Irwin",
	  "gender": "female",
	  "company": "VURBO",
	  "email": "sueirwin@vurbo.com",
	  "phone": "+1 (829) 544-2803",
	  "address": "196 Pierrepont Street, Bartonsville, Idaho, 2203",
	  "about": "Duis nisi Lorem occaecat do eu fugiat consectetur. Reprehenderit ut magna velit est reprehenderit Lorem. Excepteur consequat velit enim veniam quis velit velit enim aliquip nisi commodo ex. Pariatur consequat laboris amet fugiat nulla quis duis irure proident duis. Elit irure officia consequat reprehenderit commodo ad. Reprehenderit amet pariatur voluptate laboris dolor et veniam non ex. Consectetur ipsum in sunt irure cupidatat voluptate id ipsum.\r\n",
	  "registered": "2024-01-07T08:57:16 -09:00",
	  "latitude": 48.927862,
	  "longitude": -79.629704,
	  "tags": [
		"consequat",
		"ex",
		"tempor",
		"dolor",
		"nisi",
		"occaecat",
		"quis"
	  ],
	  "friends": [
		{
		  "id": 0,
		  "name": "Michael Rojas"
		},
		{
		  "id": 1,
		  "name": "Rodgers Pennington"
		},
		{
		  "id": 2,
		  "name": "Cardenas Monroe"
		}
	  ],
	  "greeting": "Hello, Sue Irwin! You have 7 unread messages.",
	  "favoriteFruit": "banana"
	},
	{
	  "_id": "672d31b2240100375952cb1e",
	  "index": 1,
	  "guid": "98329871-09ba-4d38-99e3-8d808395bdfa",
	  "isActive": false,
	  "balance": "$2,832.86",
	  "picture": "http://placehold.it/32x32",
	  "age": 27,
	  "eyeColor": "blue",
	  "name": "Noemi Hays",
	  "gender": "female",
	  "company": "OBLIQ",
	  "email": "noemihays@obliq.com",
	  "phone": "+1 (940) 540-2100",
	  "address": "368 Suydam Place, Ezel, Ohio, 7763",
	  "about": "Pariatur dolore commodo ex aliqua tempor qui sit. Dolor incididunt nulla anim occaecat excepteur consectetur commodo officia voluptate tempor voluptate eiusmod officia ut. Officia labore aliquip fugiat amet aliquip excepteur et qui et laboris aliquip sunt occaecat. Dolore sint amet Lorem ea. Ipsum magna esse amet culpa sunt est nostrud id ut nostrud pariatur eu sit anim. Occaecat Lorem aute elit reprehenderit est reprehenderit velit exercitation qui amet. Enim consequat incididunt est velit ad.\r\n",
	  "registered": "2023-03-12T06:12:32 -09:00",
	  "latitude": 44.210044,
	  "longitude": -165.285857,
	  "tags": [
		"exercitation",
		"pariatur",
		"amet",
		"qui",
		"deserunt"
	  ],
	  "friends": [
		{
		  "id": 0,
		  "name": "Christa Cole"
		},
		{
		  "id": 1,
		  "name": "Nona Knowles",
		  "age": 32,
		  "sex": "male",
		  "status": "student"
  
		},
		{
		  "id": 2,
		  "name": "Keisha Mosley"
		}
	  ],
	  "greeting": "Hello, Noemi Hays! You have 9 unread messages.",
	  "favoriteFruit": "strawberry"
	},
	{
	  "_id": "672d31b23cbd5ea5d1a839c6",
	  "index": 2,
	  "guid": "f6cf6351-3bdd-4a3f-a353-f612e2befa91",
	  "isActive": true,
	  "balance": "$3,136.55",
	  "picture": "http://placehold.it/32x32",
	  "age": 23,
	  "eyeColor": "brown",
	  "name": "Parks Buck",
	  "gender": "male",
	  "company": "ZOSIS",
	  "email": "parksbuck@zosis.com",
	  "phone": "+1 (842) 554-2877",
	  "address": "283 Kenmore Terrace, Shawmut, District Of Columbia, 1837",
	  "about": "Proident sint id fugiat minim ut mollit nostrud magna eiusmod commodo. Laborum culpa amet ea occaecat. Do tempor labore irure officia dolor magna ullamco reprehenderit tempor. Deserunt qui dolore reprehenderit amet dolor adipisicing fugiat exercitation ad voluptate aute. Culpa amet culpa ea dolor cillum exercitation pariatur sit eu cillum dolor fugiat.\r\n",
	  "registered": "2021-11-06T03:15:03 -09:00",
	  "latitude": -57.562182,
	  "longitude": 29.246979,
	  "tags": [
		"cupidatat",
		"do",
		"Lorem",
		"cupidatat",
		"quis",
		"reprehenderit",
		"do"
	  ],
	  "friends": [
		{
		  "id": 0,
		  "name": "Alma Neal"
		},
		{
		  "id": 1,
		  "name": "Trudy Crawford"
		},
		{
		  "id": 2,
		  "name": "Bernadette Whitaker"
		}
	  ],
	  "greeting": "Hello, Parks Buck! You have 10 unread messages.",
	  "favoriteFruit": "apple"
	},
	{
	  "_id": "672d31b2e474ce8b790e060e",
	  "index": 3,
	  "guid": "879c9095-f972-4519-af5e-e84878cc8ae9",
	  "isActive": false,
	  "balance": "$2,673.28",
	  "picture": "http://placehold.it/32x32",
	  "age": 29,
	  "eyeColor": "blue",
	  "name": "Hester Petersen",
	  "gender": "male",
	  "company": "ZERBINA",
	  "email": "hesterpetersen@zerbina.com",
	  "phone": "+1 (962) 504-3557",
	  "address": "181 Beard Street, Wedgewood, Northern Mariana Islands, 1108",
	  "about": "Ut et quis enim dolore mollit. Reprehenderit excepteur esse Lorem ea veniam voluptate non tempor. Deserunt qui cupidatat est sunt. Sit aliqua nulla mollit dolore cillum commodo magna. Dolor non esse aute duis commodo Lorem dolor quis sit.\r\n",
	  "registered": "2016-06-21T10:37:26 -09:00",
	  "latitude": 81.428735,
	  "longitude": 62.308805,
	  "tags": [
		"voluptate",
		"quis",
		"ea",
		"veniam",
		"sunt",
		"deserunt",
		"commodo"
	  ],
	  "friends": [
		{
		  "id": 0,
		  "name": "Lori Gutierrez"
		},
		{
		  "id": 1,
		  "name": "Nola Griffin"
		},
		{
		  "id": 2,
		  "name": "Corrine Gross"
		}
	  ],
	  "greeting": "Hello, Hester Petersen! You have 10 unread messages.",
	  "favoriteFruit": "strawberry"
	},
	{
	  "_id": "672d31b2aacc9b71b972f01f",
	  "index": 4,
	  "guid": "33fadf23-b3d6-4710-8be5-07d6804ca411",
	  "isActive": false,
	  "balance": "$2,411.09",
	  "picture": "http://placehold.it/32x32",
	  "age": 32,
	  "eyeColor": "brown",
	  "name": "Margret Lloyd",
	  "gender": "female",
	  "company": "ORBAXTER",
	  "email": "margretlloyd@orbaxter.com",
	  "phone": "+1 (884) 556-3656",
	  "address": "977 Schaefer Street, Olney, Connecticut, 5275",
	  "about": "Minim irure nulla duis mollit exercitation nisi aliqua tempor consectetur magna. Occaecat velit do veniam ipsum magna veniam elit. Voluptate nisi Lorem ullamco et. Voluptate elit consectetur enim Lorem in ipsum quis. In Lorem exercitation ullamco nulla magna laboris fugiat labore amet. Do mollit labore cillum est veniam amet mollit reprehenderit. Laborum enim ullamco ad ea consectetur ex occaecat ut dolore.\r\n",
	  "registered": "2021-09-28T05:31:32 -09:00",
	  "latitude": -80.987436,
	  "longitude": -126.500677,
	  "tags": [
		"ea",
		"quis",
		"aliquip",
		"aliquip",
		"sint",
		"enim",
		"quis"
	  ],
	  "friends": [
		{
		  "id": 0,
		  "name": "Daniel Pace"
		},
		{
		  "id": 1,
		  "name": "Simmons Nunez"
		},
		{
		  "id": 2,
		  "name": "Roach Manning"
		}
	  ],
	  "greeting": "Hello, Margret Lloyd! You have 2 unread messages.",
	  "favoriteFruit": "strawberry"
	},
	{
	  "_id": "672d31b224f138c8a64311fc",
	  "index": 5,
	  "guid": "2da6913b-6155-4967-988f-294cc54f64de",
	  "isActive": true,
	  "balance": "$1,152.90",
	  "picture": "http://placehold.it/32x32",
	  "age": 35,
	  "eyeColor": "green",
	  "name": "Velasquez Greer",
	  "gender": "male",
	  "company": "ACRODANCE",
	  "email": "velasquezgreer@acrodance.com",
	  "phone": "+1 (826) 408-2675",
	  "address": "588 Malbone Street, Catherine, Maine, 8474",
	  "about": "Quis incididunt consectetur pariatur ipsum deserunt ea nisi ullamco. Minim in occaecat adipisicing nisi id labore Lorem cillum cupidatat fugiat dolor commodo. Sunt pariatur ipsum dolor elit aliquip laborum veniam aliqua consectetur sit sint nisi labore.\r\n",
	  "registered": "2021-10-10T10:32:14 -09:00",
	  "latitude": -65.290247,
	  "longitude": 116.676104,
	  "tags": [
		"irure",
		"laboris",
		"consectetur",
		"deserunt",
		"aliqua",
		"deserunt",
		"fugiat"
	  ],
	  "friends": [
		{
		  "id": 0,
		  "name": "Berry Zamora"
		},
		{
		  "id": 1,
		  "name": "Viola Myers"
		},
		{
		  "id": 2,
		  "name": "Leanna Pratt"
		}
	  ],
	  "greeting": "Hello, Velasquez Greer! You have 10 unread messages.",
	  "favoriteFruit": "banana"
	}
  ]`