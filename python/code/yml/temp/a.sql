-- MySQL dump 10.13  Distrib 5.5.38, for debian-linux-gnu (x86_64)
--
-- Host: localhost    Database: speediodb
-- ------------------------------------------------------
-- Server version	5.5.38-0ubuntu0.14.04.1

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `cloudsetting`
--

DROP TABLE IF EXISTS `cloudsetting`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cloudsetting` (
  `uid` int(10) NOT NULL AUTO_INCREMENT,
  `Mysql` varchar(64) DEFAULT NULL,
  `Mongo` varchar(64) DEFAULT NULL,
  `Master` varchar(64) DEFAULT NULL,
  `Worker` varchar(64) DEFAULT NULL,
  `Service` varchar(64) DEFAULT NULL,
  `Cloudstor` varchar(64) DEFAULT NULL,
  `Store` varchar(64) DEFAULT NULL,
  PRIMARY KEY (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cloudsetting`
--

LOCK TABLES `cloudsetting` WRITE;
/*!40000 ALTER TABLE `cloudsetting` DISABLE KEYS */;
/*!40000 ALTER TABLE `cloudsetting` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `disk`
--

DROP TABLE IF EXISTS `disk`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `disk` (
  `uuid` varchar(64) DEFAULT NULL,
  `location` varchar(64) DEFAULT NULL,
  `machineId` varchar(64) DEFAULT NULL,
  `health` varchar(64) DEFAULT NULL,
  `role` varchar(64) DEFAULT NULL,
  `cap_sector` bigint(20) NOT NULL,
  `raid` varchar(64) DEFAULT NULL,
  `vendor` varchar(64) DEFAULT NULL,
  `model` varchar(64) DEFAULT NULL,
  `sn` varchar(64) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `disk`
--

LOCK TABLES `disk` WRITE;
/*!40000 ALTER TABLE `disk` DISABLE KEYS */;
INSERT INTO `disk` VALUES ('3bfc894a-74c6-43e0-af4d-860c321c40bc','1.1.8','73b9f5ca-4d73-ded8-0cb8-7c0478375aae192168276','normal','data',7814037168,'ed8d9b1d-617d-4199-869f-add83fca6b97','ST4000VX000-1F4168','UNKOWN',''),('51e45712-a5f3-4d0b-befd-1f5aca7f82af','1.1.3','73b9f5ca-4d73-ded8-0cb8-7c0478375aae192168276','normal','data',7814037168,'257013e2-a4cb-4a7e-b0d5-562f07f8972c','ST4000VX000-1F4168','UNKOWN',''),('62e7e3f4-0763-4776-8081-d6396b304ebb','1.1.2','73b9f5ca-4d73-ded8-0cb8-7c0478375aae192168276','normal','data',7814037168,'257013e2-a4cb-4a7e-b0d5-562f07f8972c','ST4000VX000-1F4168','UNKOWN',''),('6e353a5f-4c65-422e-bdb4-4590113beb55','1.1.1','73b9f5ca-4d73-ded8-0cb8-7c0478375aae192168276','normal','data',7814037168,'257013e2-a4cb-4a7e-b0d5-562f07f8972c','ST4000VX000-1F4168','UNKOWN',''),('84786c71-ca5a-4d93-bafc-0ffcc61afc70','1.1.11','73b9f5ca-4d73-ded8-0cb8-7c0478375aae192168276','normal','unused',7814037168,'','ST4000VX000-1F4168','UNKOWN',''),('9cdc1d93-04d8-486e-bb6d-84ea07bd0198','1.1.4','73b9f5ca-4d73-ded8-0cb8-7c0478375aae192168276','normal','unused',1953525168,'','WDC','WD1002F9YZ-09H1JL1','WD-WMC5K0D7D2ZS'),('cd1690a4-ecad-4893-9f57-f75a0dd69e3a','1.1.6','73b9f5ca-4d73-ded8-0cb8-7c0478375aae192168276','normal','data',7814037168,'ed8d9b1d-617d-4199-869f-add83fca6b97','ST4000VX000-1F4168','UNKOWN',''),('dd653f76-d33a-47fd-b622-2e5a9a394075','1.1.12','73b9f5ca-4d73-ded8-0cb8-7c0478375aae192168276','normal','unused',7814037168,'','ST4000VX000-1F4168','UNKOWN',''),('f4725a0b-de5f-4bea-83c9-4a9598972619','1.1.7','73b9f5ca-4d73-ded8-0cb8-7c0478375aae192168276','normal','data',7814037168,'ed8d9b1d-617d-4199-869f-add83fca6b97','ST4000VX000-1F4168','UNKOWN','');
/*!40000 ALTER TABLE `disk` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `disks`
--

DROP TABLE IF EXISTS `disks`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `disks` (
  `uuid` varchar(255) NOT NULL,
  `health` varchar(255) NOT NULL DEFAULT '',
  `role` varchar(255) NOT NULL DEFAULT '',
  `location` varchar(255) NOT NULL DEFAULT '',
  `raid` varchar(255) NOT NULL DEFAULT '',
  `cap__sector` bigint(20) NOT NULL DEFAULT '0',
  `vendor` varchar(255) NOT NULL DEFAULT '',
  `model` varchar(255) NOT NULL DEFAULT '',
  `sn` varchar(255) NOT NULL DEFAULT '',
  `cap_sector` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `disks`
--

LOCK TABLES `disks` WRITE;
/*!40000 ALTER TABLE `disks` DISABLE KEYS */;
/*!40000 ALTER TABLE `disks` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `filesystems`
--

DROP TABLE IF EXISTS `filesystems`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `filesystems` (
  `uuid` varchar(64) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `machineId` varchar(64) DEFAULT NULL,
  `volume` varchar(64) DEFAULT NULL,
  `name` varchar(64) DEFAULT NULL,
  `chunk_kb` int(11) DEFAULT NULL,
  `mountpoint` varchar(64) DEFAULT NULL,
  `type` varchar(64) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `filesystems`
--

LOCK TABLES `filesystems` WRITE;
/*!40000 ALTER TABLE `filesystems` DISABLE KEYS */;
INSERT INTO `filesystems` VALUES ('827ba09e-0237-4bfc-b933-87d04ada8443',NULL,NULL,'73b9f5ca-4d73-ded8-0cb8-7c0478375aae192168276','mfpsc3vb-ycxt-mhqi-cmji-fapyi8sw6z4t','myfs2',256,NULL,'xfs'),('82ca70da-1c63-486f-94b4-e5052f34786b',NULL,NULL,'73b9f5ca-4d73-ded8-0cb8-7c0478375aae192168276','l25huitg-zkqc-yjey-f1l8-rmtgprpehhrm','myfs3',256,NULL,'xfs'),('be9a9ca6-3dd3-4c84-bd81-b1ab897f34d9',NULL,NULL,'73b9f5ca-4d73-ded8-0cb8-7c0478375aae192168276','xsqawsmr-ja6f-hptl-1xld-p51ounxxtsj1','myfs1',256,NULL,'xfs');
/*!40000 ALTER TABLE `filesystems` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `initiator`
--

DROP TABLE IF EXISTS `initiator`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `initiator` (
  `wwn` varchar(64) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `target_wwn` varchar(64) DEFAULT NULL,
  `machineId` varchar(64) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `initiator`
--

LOCK TABLES `initiator` WRITE;
/*!40000 ALTER TABLE `initiator` DISABLE KEYS */;
INSERT INTO `initiator` VALUES ('iqn.2013-01.net.zbx.initiator:cc22da07fae7',NULL,NULL,'iqn.2013-01.net.zbx.target:cc22da07fae7a7','73b9f5ca-4d73-ded8-0cb8-7c0478375aae1921682135');
/*!40000 ALTER TABLE `initiator` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `initiator_volume`
--

DROP TABLE IF EXISTS `initiator_volume`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `initiator_volume` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `initiator` varchar(255) DEFAULT NULL,
  `volume` varchar(255) DEFAULT NULL,
  `machineId` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `initiator_volume`
--

LOCK TABLES `initiator_volume` WRITE;
/*!40000 ALTER TABLE `initiator_volume` DISABLE KEYS */;
INSERT INTO `initiator_volume` VALUES (1,'0000-00-00 00:00:00','0000-00-00 00:00:00','iqn.2013-01.net.zbx.initiator:cc22da07fae7','wxambwrh-naif-ksay-bxke-48lenj4nr2hz','73b9f5ca-4d73-ded8-0cb8-7c0478375aae1921682135');
/*!40000 ALTER TABLE `initiator_volume` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `journals`
--

DROP TABLE IF EXISTS `journals`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `journals` (
  `uid` int(10) NOT NULL AUTO_INCREMENT,
  `level` varchar(50) DEFAULT NULL,
  `message` varchar(200) DEFAULT NULL,
  `chinese_message` varchar(200) DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  PRIMARY KEY (`uid`)
) ENGINE=InnoDB AUTO_INCREMENT=58 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `journals`
--

LOCK TABLES `journals` WRITE;
/*!40000 ALTER TABLE `journals` DISABLE KEYS */;
INSERT INTO `journals` VALUES (52,'warning','Server 192.168.2.149 unset mysql %!(EXTRA string=successfully)','服务器 192.168.2.149 解除 mysql 成功','2016-10-09 03:51:14','2016-10-09 03:51:14'),(53,'info','Server 192.168.2.149 check mysql %!(EXTRA string=unsuccessfully)','服务器 192.168.2.149 检查 mysql 失败','2016-10-09 03:51:33','2016-10-09 03:51:33'),(54,'info','Server 192.168.2.149 set mysql %!(EXTRA string=successfully)','服务器 192.168.2.149 配置 mysql 成功','2016-10-09 03:51:41','2016-10-09 03:51:41'),(55,'warning','Server 192.168.2.149 unset mysql %!(EXTRA string=successfully)','服务器 192.168.2.149 解除 mysql 成功','2016-10-09 03:54:41','2016-10-09 03:54:41'),(56,'info','Server 192.168.2.149 set mysql %!(EXTRA string=successfully)','服务器 192.168.2.149 配置 mysql 成功','2016-10-09 03:55:15','2016-10-09 03:55:15'),(57,'info','Server 192.168.2.149 set mongo %!(EXTRA string=successfully)','服务器 192.168.2.149 配置 mongo 成功','2016-10-09 04:00:53','2016-10-09 04:00:53');
/*!40000 ALTER TABLE `journals` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `machine`
--

DROP TABLE IF EXISTS `machine`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `machine` (
  `uid` int(10) NOT NULL AUTO_INCREMENT,
  `uuid` varchar(64) DEFAULT NULL,
  `ip` varchar(64) DEFAULT NULL,
  `slotnr` int(10) DEFAULT NULL,
  `created` datetime DEFAULT NULL,
  PRIMARY KEY (`uid`),
  KEY `machine_created` (`created`)
) ENGINE=InnoDB AUTO_INCREMENT=35 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `machine`
--

LOCK TABLES `machine` WRITE;
/*!40000 ALTER TABLE `machine` DISABLE KEYS */;
INSERT INTO `machine` VALUES (31,'73b9f5ca-4d73-ded8-0cb8-7c0478375aae1921682102','192.168.2.102',24,'2016-10-11 02:58:45'),(32,'73b9f5ca-4d73-ded8-0cb8-7c0478375aae1921682103','192.168.2.103',24,'2016-10-11 05:53:25'),(33,'73b9f5ca-4d73-ded8-0cb8-7c0478375aae1921682135','192.168.2.135',24,'2016-10-11 07:05:14'),(34,'73b9f5ca-4d73-ded8-0cb8-7c0478375aae192168276','192.168.2.76',24,'2016-10-11 07:08:32');
/*!40000 ALTER TABLE `machine` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `network_initiator`
--

DROP TABLE IF EXISTS `network_initiator`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `network_initiator` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `initiator` varchar(255) DEFAULT NULL,
  `eth` varchar(255) DEFAULT NULL,
  `port` int(11) DEFAULT NULL,
  `machineId` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `network_initiator`
--

LOCK TABLES `network_initiator` WRITE;
/*!40000 ALTER TABLE `network_initiator` DISABLE KEYS */;
INSERT INTO `network_initiator` VALUES (1,'0000-00-00 00:00:00','0000-00-00 00:00:00','iqn.2013-01.net.zbx.initiator:cb814853fa22','eth0',3260,'73b9f5ca-4d73-ded8-0cb8-7c0478375aae1921682102');
/*!40000 ALTER TABLE `network_initiator` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `raid`
--

DROP TABLE IF EXISTS `raid`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `raid` (
  `uuid` varchar(64) DEFAULT NULL,
  `health` varchar(64) DEFAULT NULL,
  `machineId` varchar(64) DEFAULT NULL,
  `level` varchar(64) DEFAULT NULL,
  `name` varchar(64) DEFAULT NULL,
  `cap` int(11) DEFAULT NULL,
  `used_cap` int(11) DEFAULT NULL,
  `deleted` tinyint(1) DEFAULT '0'
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `raid`
--

LOCK TABLES `raid` WRITE;
/*!40000 ALTER TABLE `raid` DISABLE KEYS */;
INSERT INTO `raid` VALUES ('24c8ac86-3b98-4d08-9ad8-4e00c5bc3895','normal','73b9f5ca-4d73-ded8-0cb8-7c0478375aae192168276','0','sddd',3725,0,1),('257013e2-a4cb-4a7e-b0d5-562f07f8972c','normal','73b9f5ca-4d73-ded8-0cb8-7c0478375aae192168276','5','rd8',7451,7451,0),('ed8d9b1d-617d-4199-869f-add83fca6b97','normal','73b9f5ca-4d73-ded8-0cb8-7c0478375aae192168276','5','rd5',7451,6000,0);
/*!40000 ALTER TABLE `raid` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `raid_volume`
--

DROP TABLE IF EXISTS `raid_volume`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `raid_volume` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `raid` varchar(255) DEFAULT NULL,
  `volume` varchar(255) DEFAULT NULL,
  `type` varchar(255) NOT NULL,
  `machineId` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `raid_volume`
--

LOCK TABLES `raid_volume` WRITE;
/*!40000 ALTER TABLE `raid_volume` DISABLE KEYS */;
INSERT INTO `raid_volume` VALUES (1,'0000-00-00 00:00:00','0000-00-00 00:00:00','9e9125c9-49b3-406d-bb92-3e51d5dee848','ckw2yo2s-ky7g-hfqc-jla3-mizznylrlmpb','normal','73b9f5ca-4d73-ded8-0cb8-7c0478375aae1921682103'),(2,'0000-00-00 00:00:00','0000-00-00 00:00:00','bd3140e7-739d-4944-892b-9d161b9828d9','ckecwk46-oyty-6cvk-iebs-jb66xjmm1ry3','normal','73b9f5ca-4d73-ded8-0cb8-7c0478375aae1921682103');
/*!40000 ALTER TABLE `raid_volume` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `setting`
--

DROP TABLE IF EXISTS `setting`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `setting` (
  `uid` int(10) NOT NULL AUTO_INCREMENT,
  `Settingtype` varchar(64) DEFAULT NULL,
  `Ip` varchar(64) DEFAULT NULL,
  `Status` varchar(64) DEFAULT NULL,
  PRIMARY KEY (`uid`)
) ENGINE=InnoDB AUTO_INCREMENT=23 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `setting`
--

LOCK TABLES `setting` WRITE;
/*!40000 ALTER TABLE `setting` DISABLE KEYS */;
/*!40000 ALTER TABLE `setting` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `volume`
--

DROP TABLE IF EXISTS `volume`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `volume` (
  `uuid` varchar(64) DEFAULT NULL,
  `health` varchar(64) DEFAULT NULL,
  `machineId` varchar(64) DEFAULT NULL,
  `name` varchar(64) DEFAULT NULL,
  `used` int(11) DEFAULT NULL,
  `owner_type` varchar(64) DEFAULT NULL,
  `cap` bigint(20) NOT NULL DEFAULT '0',
  `deleted` tinyint(1) DEFAULT '0'
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `volume`
--

LOCK TABLES `volume` WRITE;
/*!40000 ALTER TABLE `volume` DISABLE KEYS */;
INSERT INTO `volume` VALUES ('bv7qj83o-hxid-ceov-wus8-awol4v6bhnzy','normal','73b9f5ca-4d73-ded8-0cb8-7c0478375aae192168276','asd',0,'',145,2),('gfgdrlgg-vdnv-tsoe-n8p9-57rmfbpxsgfh','normal','73b9f5ca-4d73-ded8-0cb8-7c0478375aae192168276','asd',0,'',145,2),('gsbscf2r-ite2-pzkl-ydrl-rbq90iwyg8js','normal','73b9f5ca-4d73-ded8-0cb8-7c0478375aae192168276','sda',0,'',145,2),('l25huitg-zkqc-yjey-f1l8-rmtgprpehhrm','normal','73b9f5ca-4d73-ded8-0cb8-7c0478375aae192168276','vol2',1,'xfs',3000,0),('l4kkmwr0-ta6c-hteo-ovj3-zn5lbfo4fdv3','normal','73b9f5ca-4d73-ded8-0cb8-7c0478375aae192168276','sd',0,'',14,2),('lxamwop2-brzg-jr93-vibt-fcklcb1cuv8n','normal','73b9f5ca-4d73-ded8-0cb8-7c0478375aae192168276','vv',0,'',145,2),('mfpsc3vb-ycxt-mhqi-cmji-fapyi8sw6z4t','normal','73b9f5ca-4d73-ded8-0cb8-7c0478375aae192168276','vol0',1,'xfs',3000,0),('n1c6xfdl-bot7-zi30-rpco-tmtlti2pmtg4','normal','73b9f5ca-4d73-ded8-0cb8-7c0478375aae192168276','cava',0,'',145,2),('q0r3sh36-0rjw-9lur-nbaf-tgquuarcjh4e','normal','73b9f5ca-4d73-ded8-0cb8-7c0478375aae192168276','fas',0,'',145,2),('xet0zffi-opfn-zs6b-df43-tqlbxbbf988p','normal','73b9f5ca-4d73-ded8-0cb8-7c0478375aae192168276','dsa',0,'',130,2),('xsqawsmr-ja6f-hptl-1xld-p51ounxxtsj1','normal','73b9f5ca-4d73-ded8-0cb8-7c0478375aae192168276','vol8',1,'xfs',7451,0),('yoxkbjwa-jocd-a0nu-ljsq-zxvpijiafh0d','normal','73b9f5ca-4d73-ded8-0cb8-7c0478375aae192168276','sa',0,'',14,2),('z1l3quji-uelk-y04h-jdsv-0ftkoge1qek7','normal','73b9f5ca-4d73-ded8-0cb8-7c0478375aae192168276','asdfa',0,'',145,2);
/*!40000 ALTER TABLE `volume` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2016-10-12 19:57:57
