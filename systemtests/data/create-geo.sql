# Copyright 2026 Richard Kosegi
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


CREATE TABLE continent (
    id INT(10) UNSIGNED NOT NULL AUTO_INCREMENT ,
    name VARCHAR(200) NOT NULL ,
    size INT UNSIGNED NOT NULL ,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
;

CREATE TABLE country (
    id INT(10) UNSIGNED NOT NULL AUTO_INCREMENT ,
    name VARCHAR(200) NOT NULL ,
    size INT UNSIGNED NOT NULL ,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
;

CREATE TABLE city (
    id INT(10) UNSIGNED NOT NULL AUTO_INCREMENT ,
    country_id INT(10) UNSIGNED NOT NULL ,
    name VARCHAR(200) NOT NULL ,
    PRIMARY KEY (`id`) ,
    CONSTRAINT `city_country_fk` FOREIGN KEY (`country_id`) REFERENCES `country` (`id`)
) ENGINE = InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
;

-- continents
INSERT INTO continent(id,name, size) VALUES(1, 'Asia', 44614000);
INSERT INTO continent(id,name, size) VALUES(2, 'Africa', 30365000);
INSERT INTO continent(id,name, size) VALUES(3, 'North America', 24230000);
INSERT INTO continent(id,name, size) VALUES(4, 'South America', 17814000);
INSERT INTO continent(id,name, size) VALUES(5, 'Antarctica', 14200000);
INSERT INTO continent(id,name, size) VALUES(6, 'Europe', 10000000);
INSERT INTO continent(id,name, size) VALUES(7, 'Oceania', 8510926);

-- countries
INSERT INTO country(id,name,size) VALUES(1,'Russia',17098246);
INSERT INTO country(id,name,size) VALUES(2,'Canada',9984670);
INSERT INTO country(id,name,size) VALUES(3,'China',9596960);
INSERT INTO country(id,name,size) VALUES(4,'United States',9525067);
INSERT INTO country(id,name,size) VALUES(5,'Brazil',8510346);
INSERT INTO country(id,name,size) VALUES(6,'Australia',7741220);
INSERT INTO country(id,name,size) VALUES(7,'India',3287263);
INSERT INTO country(id,name,size) VALUES(8,'Argentina',2780400);
INSERT INTO country(id,name,size) VALUES(9,'Kazakhstan',2724910);
INSERT INTO country(id,name,size) VALUES(10,'Algeria',2381741);
INSERT INTO country(id,name,size) VALUES(11,'Democratic Republic of the Congo',2344858);
INSERT INTO country(id,name,size) VALUES(12,'Saudi Arabia',2149690);
INSERT INTO country(id,name,size) VALUES(13,'Mexico',1964375);
INSERT INTO country(id,name,size) VALUES(14,'Indonesia',1904569);
INSERT INTO country(id,name,size) VALUES(15,'Sudan',1878000);
INSERT INTO country(id,name,size) VALUES(16,'Libya',1759540);
INSERT INTO country(id,name,size) VALUES(17,'Iran',1648195);
INSERT INTO country(id,name,size) VALUES(18,'Mongolia',1564116);
INSERT INTO country(id,name,size) VALUES(19,'Peru',1285216);
INSERT INTO country(id,name,size) VALUES(20,'Chad',1284000);
INSERT INTO country(id,name,size) VALUES(21,'Niger',1267000);
INSERT INTO country(id,name,size) VALUES(22,'Angola',1246700);
INSERT INTO country(id,name,size) VALUES(23,'Mali',1240192);
INSERT INTO country(id,name,size) VALUES(24,'South Africa',1219090);
INSERT INTO country(id,name,size) VALUES(25,'Colombia',1138910);
INSERT INTO country(id,name,size) VALUES(26,'Ethiopia',1104300);
INSERT INTO country(id,name,size) VALUES(27,'Bolivia',1098581);
INSERT INTO country(id,name,size) VALUES(28,'Mauritania',1030700);
INSERT INTO country(id,name,size) VALUES(29,'Egypt',1001450);
INSERT INTO country(id,name,size) VALUES(30,'Tanzania',947303);
INSERT INTO country(id,name,size) VALUES(31,'Nigeria',923768);
INSERT INTO country(id,name,size) VALUES(32,'Venezuela',912050);
INSERT INTO country(id,name,size) VALUES(33,'Pakistan',882363);
INSERT INTO country(id,name,size) VALUES(34,'Namibia',824292);
INSERT INTO country(id,name,size) VALUES(35,'Mozambique',799380);
INSERT INTO country(id,name,size) VALUES(36,'Turkey',783562);
INSERT INTO country(id,name,size) VALUES(37,'Chile',756102);
INSERT INTO country(id,name,size) VALUES(38,'Zambia',752612);
INSERT INTO country(id,name,size) VALUES(39,'Myanmar',676578);
INSERT INTO country(id,name,size) VALUES(40,'Afghanistan',652864);
INSERT INTO country(id,name,size) VALUES(41,'South Sudan',644329);
INSERT INTO country(id,name,size) VALUES(42,'France',643801);
INSERT INTO country(id,name,size) VALUES(43,'Somalia',637657);
INSERT INTO country(id,name,size) VALUES(44,'Central African Republic',622984);
INSERT INTO country(id,name,size) VALUES(45,'Ukraine',603550);
INSERT INTO country(id,name,size) VALUES(46,'Madagascar',587041);
INSERT INTO country(id,name,size) VALUES(47,'Botswana',582000);
INSERT INTO country(id,name,size) VALUES(48,'Kenya',580367);
INSERT INTO country(id,name,size) VALUES(49,'Thailand',513120);
INSERT INTO country(id,name,size) VALUES(50,'Spain',505370);
INSERT INTO country(id,name,size) VALUES(51,'Turkmenistan',488100);
INSERT INTO country(id,name,size) VALUES(52,'Cameroon',475650);
INSERT INTO country(id,name,size) VALUES(53,'Papua New Guinea',462840);
INSERT INTO country(id,name,size) VALUES(54,'Yemen',455503);
INSERT INTO country(id,name,size) VALUES(55,'Sweden',450295);
INSERT INTO country(id,name,size) VALUES(56,'Uzbekistan',447400);
INSERT INTO country(id,name,size) VALUES(57,'Morocco',446550);
INSERT INTO country(id,name,size) VALUES(58,'Iraq',438317);
INSERT INTO country(id,name,size) VALUES(59,'Paraguay',406752);
INSERT INTO country(id,name,size) VALUES(60,'Zimbabwe',390757);
INSERT INTO country(id,name,size) VALUES(61,'Norway',386224);
INSERT INTO country(id,name,size) VALUES(62,'Japan',377915);
INSERT INTO country(id,name,size) VALUES(63,'Germany',357581);
INSERT INTO country(id,name,size) VALUES(64,'Republic of the Congo',342000);
INSERT INTO country(id,name,size) VALUES(65,'Finland',338145);
INSERT INTO country(id,name,size) VALUES(66,'Vietnam',331340);
INSERT INTO country(id,name,size) VALUES(67,'Malaysia',330621);
INSERT INTO country(id,name,size) VALUES(68,'Ivory Coast',322462);
INSERT INTO country(id,name,size) VALUES(69,'Poland',312685);
INSERT INTO country(id,name,size) VALUES(70,'Oman',309500);
INSERT INTO country(id,name,size) VALUES(71,'Italy',302068);
INSERT INTO country(id,name,size) VALUES(72,'Philippines',300000);
INSERT INTO country(id,name,size) VALUES(73,'Ecuador',283561);
INSERT INTO country(id,name,size) VALUES(74,'Burkina Faso',274200);
INSERT INTO country(id,name,size) VALUES(75,'New Zealand',268838);
INSERT INTO country(id,name,size) VALUES(76,'Gabon',267668);
INSERT INTO country(id,name,size) VALUES(77,'Guinea',245857);
INSERT INTO country(id,name,size) VALUES(78,'United Kingdom',244376);
INSERT INTO country(id,name,size) VALUES(79,'Uganda',241550);
INSERT INTO country(id,name,size) VALUES(80,'Ghana',238537);
INSERT INTO country(id,name,size) VALUES(81,'Romania',238398);
INSERT INTO country(id,name,size) VALUES(82,'Laos',236800);
INSERT INTO country(id,name,size) VALUES(83,'Guyana',214969);
INSERT INTO country(id,name,size) VALUES(84,'Belarus',207600);
INSERT INTO country(id,name,size) VALUES(85,'Kyrgyzstan',199949);
INSERT INTO country(id,name,size) VALUES(86,'Senegal',196712);
INSERT INTO country(id,name,size) VALUES(87,'Syria',185180);
INSERT INTO country(id,name,size) VALUES(88,'Cambodia',181035);
INSERT INTO country(id,name,size) VALUES(89,'Uruguay',176215);
INSERT INTO country(id,name,size) VALUES(90,'Suriname',163820);
INSERT INTO country(id,name,size) VALUES(91,'Tunisia',163610);
INSERT INTO country(id,name,size) VALUES(92,'Bangladesh',148460);
INSERT INTO country(id,name,size) VALUES(93,'Nepal',147181);
INSERT INTO country(id,name,size) VALUES(94,'Tajikistan',144100);
INSERT INTO country(id,name,size) VALUES(95,'Greece',131957);
INSERT INTO country(id,name,size) VALUES(96,'Nicaragua',130373);
INSERT INTO country(id,name,size) VALUES(97,'North Korea',120538);
INSERT INTO country(id,name,size) VALUES(98,'Malawi',118484);
INSERT INTO country(id,name,size) VALUES(99,'Eritrea',117600);
INSERT INTO country(id,name,size) VALUES(100,'Benin',114763);
INSERT INTO country(id,name,size) VALUES(101,'Honduras',112492);
INSERT INTO country(id,name,size) VALUES(102,'Liberia',111369);
INSERT INTO country(id,name,size) VALUES(103,'Bulgaria',110879);
INSERT INTO country(id,name,size) VALUES(104,'Cuba',109884);
INSERT INTO country(id,name,size) VALUES(105,'Guatemala',108889);
INSERT INTO country(id,name,size) VALUES(106,'Iceland',103000);
INSERT INTO country(id,name,size) VALUES(107,'South Korea',100432);
INSERT INTO country(id,name,size) VALUES(108,'Hungary',93025);
INSERT INTO country(id,name,size) VALUES(109,'Portugal',92090);
INSERT INTO country(id,name,size) VALUES(110,'Jordan',89318);
INSERT INTO country(id,name,size) VALUES(111,'Serbia',88499);
INSERT INTO country(id,name,size) VALUES(112,'Azerbaijan',86600);
INSERT INTO country(id,name,size) VALUES(113,'Austria',83878);
INSERT INTO country(id,name,size) VALUES(114,'United Arab Emirates',83600);
INSERT INTO country(id,name,size) VALUES(115,'Czech Republic',78871);
INSERT INTO country(id,name,size) VALUES(116,'Panama',75320);
INSERT INTO country(id,name,size) VALUES(117,'Sierra Leone',72300);
INSERT INTO country(id,name,size) VALUES(118,'Ireland',70273);
INSERT INTO country(id,name,size) VALUES(119,'Georgia',69700);
INSERT INTO country(id,name,size) VALUES(120,'Sri Lanka',67240);
INSERT INTO country(id,name,size) VALUES(121,'Lithuania',65286);
INSERT INTO country(id,name,size) VALUES(122,'Latvia',64594);
INSERT INTO country(id,name,size) VALUES(123,'Togo',56785);
INSERT INTO country(id,name,size) VALUES(124,'Croatia',56594);
INSERT INTO country(id,name,size) VALUES(125,'Bosnia and Herzegovina',51209);
INSERT INTO country(id,name,size) VALUES(126,'Costa Rica',51180);
INSERT INTO country(id,name,size) VALUES(127,'Slovakia',49035);
INSERT INTO country(id,name,size) VALUES(128,'Dominican Republic',48670);
INSERT INTO country(id,name,size) VALUES(129,'Estonia',45339);
INSERT INTO country(id,name,size) VALUES(130,'Denmark',42947);
INSERT INTO country(id,name,size) VALUES(131,'Netherlands',41865);
INSERT INTO country(id,name,size) VALUES(132,'Switzerland',41291);
INSERT INTO country(id,name,size) VALUES(133,'Bhutan',38394);
INSERT INTO country(id,name,size) VALUES(134,'Guinea-Bissau',36125);
INSERT INTO country(id,name,size) VALUES(135,'Moldova',33847);
INSERT INTO country(id,name,size) VALUES(136,'Belgium',30528);
INSERT INTO country(id,name,size) VALUES(137,'Lesotho',30355);
INSERT INTO country(id,name,size) VALUES(138,'Armenia',29743);
INSERT INTO country(id,name,size) VALUES(139,'Solomon Islands',28896);
INSERT INTO country(id,name,size) VALUES(140,'Albania',28748);
INSERT INTO country(id,name,size) VALUES(141,'Equatorial Guinea',28051);
INSERT INTO country(id,name,size) VALUES(142,'Burundi',27834);
INSERT INTO country(id,name,size) VALUES(143,'Haiti',27750);
INSERT INTO country(id,name,size) VALUES(144,'Rwanda',26338);
INSERT INTO country(id,name,size) VALUES(145,'North Macedonia',25713);
INSERT INTO country(id,name,size) VALUES(146,'Djibouti',23200);
INSERT INTO country(id,name,size) VALUES(147,'Belize',22965);
INSERT INTO country(id,name,size) VALUES(148,'Israel',21937);
INSERT INTO country(id,name,size) VALUES(149,'El Salvador',21041);
INSERT INTO country(id,name,size) VALUES(150,'Slovenia',20273);
INSERT INTO country(id,name,size) VALUES(151,'Fiji',18272);
INSERT INTO country(id,name,size) VALUES(152,'Kuwait',17818);
INSERT INTO country(id,name,size) VALUES(153,'Eswatini',17363);
INSERT INTO country(id,name,size) VALUES(154,'Timor-Leste',14874);
INSERT INTO country(id,name,size) VALUES(155,'Montenegro',13888);
INSERT INTO country(id,name,size) VALUES(156,'Bahamas',13880);
INSERT INTO country(id,name,size) VALUES(157,'Vanuatu',12189);
INSERT INTO country(id,name,size) VALUES(158,'Qatar',11586);
INSERT INTO country(id,name,size) VALUES(159,'The Gambia',11295);
INSERT INTO country(id,name,size) VALUES(160,'Jamaica',10991);
INSERT INTO country(id,name,size) VALUES(161,'Lebanon',10452);
INSERT INTO country(id,name,size) VALUES(162,'Cyprus',9251);
INSERT INTO country(id,name,size) VALUES(163,'Palestine',6020);
INSERT INTO country(id,name,size) VALUES(164,'Brunei',5765);
INSERT INTO country(id,name,size) VALUES(165,'Trinidad and Tobago',5127);
INSERT INTO country(id,name,size) VALUES(166,'Cape Verde',4033);
INSERT INTO country(id,name,size) VALUES(167,'Samoa',2842);
INSERT INTO country(id,name,size) VALUES(168,'Luxembourg',2586);
INSERT INTO country(id,name,size) VALUES(169,'Mauritius',2096);
INSERT INTO country(id,name,size) VALUES(170,'Comoros',1861);
INSERT INTO country(id,name,size) VALUES(171,'São Tomé and Príncipe',964);
INSERT INTO country(id,name,size) VALUES(172,'Kiribati',811);
INSERT INTO country(id,name,size) VALUES(173,'Bahrain',778);
INSERT INTO country(id,name,size) VALUES(174,'Dominica',750);
INSERT INTO country(id,name,size) VALUES(175,'Tonga',747);
INSERT INTO country(id,name,size) VALUES(176,'Singapore',745);
INSERT INTO country(id,name,size) VALUES(177,'Micronesia',702);
INSERT INTO country(id,name,size) VALUES(178,'Saint Lucia',616);
INSERT INTO country(id,name,size) VALUES(179,'Andorra',468);
INSERT INTO country(id,name,size) VALUES(180,'Palau',459);
INSERT INTO country(id,name,size) VALUES(181,'Seychelles',457);
INSERT INTO country(id,name,size) VALUES(182,'Antigua and Barbuda',442);
INSERT INTO country(id,name,size) VALUES(183,'Barbados',431);
INSERT INTO country(id,name,size) VALUES(184,'Saint Vincent and the Grenadines',389);
INSERT INTO country(id,name,size) VALUES(185,'Grenada',345);
INSERT INTO country(id,name,size) VALUES(186,'Malta',315);
INSERT INTO country(id,name,size) VALUES(187,'Maldives',300);
INSERT INTO country(id,name,size) VALUES(188,'Saint Kitts and Nevis',261);
INSERT INTO country(id,name,size) VALUES(189,'Marshall Islands',181);
INSERT INTO country(id,name,size) VALUES(190,'Liechtenstein',160);
INSERT INTO country(id,name,size) VALUES(191,'San Marino',61);
INSERT INTO country(id,name,size) VALUES(192,'Tuvalu',26);
INSERT INTO country(id,name,size) VALUES(193,'Nauru',21);
INSERT INTO country(id,name,size) VALUES(194,'Monaco',2);
INSERT INTO country(id,name,size) VALUES(195,'Vatican City',0.44);
INSERT INTO country(id,name,size) VALUES(196,'Taiwan',0.44);
INSERT INTO country(id,name,size) VALUES(197,'Hong Kong',0.44);


-- cities
INSERT INTO city(id,name,country_id) VALUES(1,'Jakarta',14);
INSERT INTO city(id,name,country_id) VALUES(2,'Dhaka',92);
INSERT INTO city(id,name,country_id) VALUES(3,'Tokyo',62);
INSERT INTO city(id,name,country_id) VALUES(4,'Delhi',7);
INSERT INTO city(id,name,country_id) VALUES(5,'Shanghai',3);
INSERT INTO city(id,name,country_id) VALUES(6,'Guangzhou',3);
INSERT INTO city(id,name,country_id) VALUES(7,'Cairo',29);
INSERT INTO city(id,name,country_id) VALUES(8,'Manila',72);
INSERT INTO city(id,name,country_id) VALUES(9,'Kolkata',7);
INSERT INTO city(id,name,country_id) VALUES(10,'Seoul',107);
INSERT INTO city(id,name,country_id) VALUES(11,'Karachi',33);
INSERT INTO city(id,name,country_id) VALUES(12,'Mumbai',7);
INSERT INTO city(id,name,country_id) VALUES(13,'São Paulo',5);
INSERT INTO city(id,name,country_id) VALUES(14,'Bangkok',49);
INSERT INTO city(id,name,country_id) VALUES(15,'Mexico City',13);
INSERT INTO city(id,name,country_id) VALUES(16,'Beijing',3);
INSERT INTO city(id,name,country_id) VALUES(17,'Lahore',33);
INSERT INTO city(id,name,country_id) VALUES(18,'Istanbul',36);
INSERT INTO city(id,name,country_id) VALUES(19,'Moscow',1);
INSERT INTO city(id,name,country_id) VALUES(20,'Ho Chi Minh City',66);
INSERT INTO city(id,name,country_id) VALUES(21,'Buenos Aires',8);
INSERT INTO city(id,name,country_id) VALUES(22,'New York City',4);
INSERT INTO city(id,name,country_id) VALUES(23,'Shenzhen',3);
INSERT INTO city(id,name,country_id) VALUES(24,'Bengaluru',7);
INSERT INTO city(id,name,country_id) VALUES(25,'Osaka',62);
INSERT INTO city(id,name,country_id) VALUES(26,'Lagos',31);
INSERT INTO city(id,name,country_id) VALUES(27,'Los Angeles',4);
INSERT INTO city(id,name,country_id) VALUES(28,'Chennai',7);
INSERT INTO city(id,name,country_id) VALUES(29,'Kinshasa',11);
INSERT INTO city(id,name,country_id) VALUES(30,'Bogotá',25);
INSERT INTO city(id,name,country_id) VALUES(31,'Lima',19);
INSERT INTO city(id,name,country_id) VALUES(32,'London',78);
INSERT INTO city(id,name,country_id) VALUES(33,'Rio de Janeiro',5);
INSERT INTO city(id,name,country_id) VALUES(34,'Paris',42);
INSERT INTO city(id,name,country_id) VALUES(35,'Hyderabad',7);
INSERT INTO city(id,name,country_id) VALUES(36,'Tehran',17);
INSERT INTO city(id,name,country_id) VALUES(37,'Taipei',196);
INSERT INTO city(id,name,country_id) VALUES(38,'Luanda',22);
INSERT INTO city(id,name,country_id) VALUES(39,'Bandung',14);
INSERT INTO city(id,name,country_id) VALUES(40,'Kuala Lumpur',67);
INSERT INTO city(id,name,country_id) VALUES(41,'Dar es Salaam',30);
INSERT INTO city(id,name,country_id) VALUES(42,'Suzhou',3);
INSERT INTO city(id,name,country_id) VALUES(43,'Ahmedabad',7);
INSERT INTO city(id,name,country_id) VALUES(44,'Hangzhou',3);
INSERT INTO city(id,name,country_id) VALUES(45,'Wuhan',3);
INSERT INTO city(id,name,country_id) VALUES(46,'Tianjin',3);
INSERT INTO city(id,name,country_id) VALUES(47,'Alexandria',29);
INSERT INTO city(id,name,country_id) VALUES(48,'Nagoya',62);
INSERT INTO city(id,name,country_id) VALUES(49,'Johannesburg',24);
INSERT INTO city(id,name,country_id) VALUES(50,'Chongqing',3);
INSERT INTO city(id,name,country_id) VALUES(51,'Riyadh',12);
INSERT INTO city(id,name,country_id) VALUES(52,'Surat',7);
INSERT INTO city(id,name,country_id) VALUES(53,'Surabaya',14);
INSERT INTO city(id,name,country_id) VALUES(54,'Pune',7);
INSERT INTO city(id,name,country_id) VALUES(55,'Khartoum',15);
INSERT INTO city(id,name,country_id) VALUES(56,'Nanjing',3);
INSERT INTO city(id,name,country_id) VALUES(57,'Santiago',37);
INSERT INTO city(id,name,country_id) VALUES(58,'Chicago',4);
INSERT INTO city(id,name,country_id) VALUES(59,'Chengdu',3);
INSERT INTO city(id,name,country_id) VALUES(60,'Xian',3);
INSERT INTO city(id,name,country_id) VALUES(61,'Hong Kong',197);
INSERT INTO city(id,name,country_id) VALUES(62,'Dongguan',3);
INSERT INTO city(id,name,country_id) VALUES(63,'Foshan',3);
INSERT INTO city(id,name,country_id) VALUES(64,'Shenyang',3);
INSERT INTO city(id,name,country_id) VALUES(65,'Baghdad',58);
INSERT INTO city(id,name,country_id) VALUES(66,'Addis Ababa',26);
INSERT INTO city(id,name,country_id) VALUES(67,'Madrid',50);
INSERT INTO city(id,name,country_id) VALUES(68,'Harbin',3);
INSERT INTO city(id,name,country_id) VALUES(69,'Houston',4);
INSERT INTO city(id,name,country_id) VALUES(70,'Dallas',4);
INSERT INTO city(id,name,country_id) VALUES(71,'Toronto',2);
INSERT INTO city(id,name,country_id) VALUES(72,'Miami',4);
INSERT INTO city(id,name,country_id) VALUES(73,'Belo Horizonte',5);
INSERT INTO city(id,name,country_id) VALUES(74,'Singapore',176);
INSERT INTO city(id,name,country_id) VALUES(75,'Philadelphia',4);
INSERT INTO city(id,name,country_id) VALUES(76,'Atlanta',4);
INSERT INTO city(id,name,country_id) VALUES(77,'Fukuoka',62);
INSERT INTO city(id,name,country_id) VALUES(78,'Barcelona',50);
INSERT INTO city(id,name,country_id) VALUES(79,'Saint Petersburg',1);
INSERT INTO city(id,name,country_id) VALUES(80,'Qingdao',3);
INSERT INTO city(id,name,country_id) VALUES(81,'Dalian',3);
INSERT INTO city(id,name,country_id) VALUES(82,'Washington, D.C.',4);
INSERT INTO city(id,name,country_id) VALUES(83,'Yangon',39);
INSERT INTO city(id,name,country_id) VALUES(84,'Jinan',3);
INSERT INTO city(id,name,country_id) VALUES(85,'Guadalajara',13);
