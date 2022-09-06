package adapter

import (
	"gitlab.altex.ro/ams/go_ddwcs/awb/service/mapper"
	"gitlab.altex.ro/ams/go_ddwcs/awb/util"
	utilLib "gitlab.altex.ro/plug/go_atx_lib/util"
)

/**
 * UrgentCargus possible statuses
 *
 * 1	Adresa gresita: Companie necunoscuta la adresa de livrare
 * 2	Adresa gresita: Alta localitate
 * 3	Adresa gresita: Numar strada gresit
 * 4	Adresa gresita: Adresa incompleta
 * 5	Spre destinatar
 * 6	Adresa gresita: Lipsa persoana de contact
 * 8	Adresa gresita: Persoana necunoscuta
 * 10	Intrare in depozit
 * 11	DF Iesire din depozit catre tranzit
 * 14	Adresa gresita: Telefon gresit sau incomplet
 * 18	Adresa gresita: Strada necunoscuta
 * 20	Livrare stabilita: Livrare stabilita la data
 * 21	Livrat la destinatar (confirmat)
 * 22	Livrare stabilita: Sarbatoare legala pana la data de
 * 24	Inchis la sosirea curierului: Inchis la destinatar
 * 25	Inregistrare pe borderou de depunere CC
 * 27	Inchis la sosirea curierului: Concediu pana la
 * 32	Inchis la sosirea curierului: Program terminat
 * 34	Inchis la sosirea curierului: Inainte de inceperea programului
 * 36	Inchis la sosirea curierului: Pauza de masa
 * 38	Se asteapta ridicarea coletului de la sediu: Destinatar contactat
 * 42	Se asteapta ridicarea coletului de la sediu: Solicitare destinatar livrare la data de
 * 44	Se asteapta ridicarea coletului de la sediu: Destinatarul nu poate fi contactat
 * 50	Destinatar mutat: Companie mutata
 * 52	Destinatar mutat: Persoana mutata
 * 54	Expediere inchisa: Negasit
 * 60	Livrat deteriorat: Continut deteriorat
 * 64	Livrat deteriorat: Ambalaj deteriorat
 * 70	Ridicare din comanda client
 * 72	Ridicare cod in afara comenzii client
 * 73	Ridicare din comanda client fara scanare
 * 74	Intrare in depozit cu transfer din site
 * 75	Eroare intrare in depozit - cod inexistent
 * 76	Intrare in depozit deja confirmat
 * 77	Revenire in depozit avizat (nelivrat)
 * 79	Nelivrat: Stationat in depozit (dupa doua incercari)
 * 89	Expediere distrusa: Expediere abandonata cu acordul expeditorului
 * 93	Rutat gresit: Sortat gresit in locatie
 * 94	Rutat gresit: Incarcata gresit - origine
 * 95	Expediere nelivrata: Vreme nefavorabila
 * 96	Expediere nelivrata: La solicitarea destinatarului
 * 97	Expediere nelivrata: Ruta supra incarcata
 * 99	Expediere nelivrata: Fara acces
 * 100	Expediere nelivrata: Solicitare UrgentCargus
 * 101	Expediere nelivrata: Trafic ingreunat
 * 103	Expediere nelivrata: Lipsa act de identitate
 * 105	Expediere nelivrata: Probleme auto
 * 106	Expediere nelivrata: Calamitate naturala
 * 108	Destinatar negasit la domiciliu: Imposibil de avizat
 * 110	Destinatar negasit la domiciliu: Avizat cu bileta la adresa
 * 111	Destinatar negasit la domiciliu: Trimis SMS
 * 114	In asteptare: Se asteapta livrarea
 * 120	In asteptare: Motive tehnice
 * 121	In asteptare: Initial sortat - ruta gresita
 * 122	In asteptare: Lipsa NT
 * 124	In asteptare: Detalii destinatar incorecte
 * 127	In asteptare: Cerere destinatar/expeditor
 * 128	In asteptare: Lipsa informatii
 * 129	In asteptare: Lipsa piesa
 * 130	In asteptare: Reimpachetare
 * 132	In asteptare: Sosit cu intarziere
 * 133	In asteptare: Asteapta instructiuni
 * 136	In asteptare: Greutate/dimensiune peste limita
 * 140	In asteptare: Vreme nefavorabila
 * 141	In asteptare: Transport anulat
 * 142	In asteptare: Supraincarcat
 * 145	In asteptare: Pierdere conexiune transport
 * 146	Livrat partial: Livrat
 * 149	Intrat in depozit din cursa de tranzit
 * 153	Expediere refuzata: Refuz plata
 * 158	Expediere refuzata: Lipsa din continut
 * 159	Expediere refuzata: Deteriorat
 * 161	Expediere refuzata: Lipsa numerar (lei)
 * 180	Expediere refuzata: Expediere nesolicitata
 * 181	Expediere refuzata: Livrare intarziata
 * 182	Expediere refuzata: Livrare personala
 * 185	Expediere refuzata: Lipsa piesa
 * 186	Expediere refuzata: Colet refuzat dupa deschidere
 * 196	Expediere returnata: Solicitare expeditor
 * 197	Expediere returnata: Solicitare destinatar
 * 198	Expediere returnata: Solicitare UrgentCargus
 * 199	Livrare reprogramata: Destinatarul solicita livrare la data
 * 200	Inspectie de securitate: Inspectat de
 * 203	Transport intarziat: Vreme nefavorabila
 * 204	Transport intarziat: Trafic ingreunat
 * 206	Transport intarziat: Pierdere conexiune transport
 * 208	Transport intarziat: Supraincarcat
 * 212	Transport intarziat: Plecat cu intarziere
 * 213	Expediere livrata de o a treia parte - nu se asteapta detalii de livrare: Data estimata de livrare
 * 216	Expediere reprogramata pentru urmatorul transport : Plecare programata
 * 223	Sortare in tara de tranzit
 * 233	Intrare in depozit (info in sistem inexistent)
 * 234	PickUp Colectat (info in sistem inexistent)
 * 249	Scanat de intrare in stoc
 * 251	AR Deconsolidare in centru
 * 255	AF Sosire in centru
 * 256	Rutat gresit: Tara destinatar gresita
 */

//ignored events
var ignoredEventIds = []int{
	1,
	2,
	3,
	4,
	6,
	8,
	14,
	18,
	50,
	52,
	70,
	72,
	73,
	93,
	94,
	136,
	142,
	223,
	256,
}

var urgentCargusEventMap = map[int]int{
	38:  StatusAwbGenerated,
	42:  StatusAwbGenerated,
	44:  StatusAwbGenerated,
	130: StatusAwbGenerated,
	136: StatusAwbGenerated,

	234: StatusPickedUp,

	10:  StatusInCourierWarehouse,
	11:  StatusInCourierWarehouse,
	20:  StatusInCourierWarehouse,
	22:  StatusInCourierWarehouse,
	74:  StatusInCourierWarehouse,
	75:  StatusInCourierWarehouse,
	76:  StatusInCourierWarehouse,
	77:  StatusInCourierWarehouse,
	79:  StatusInCourierWarehouse,
	95:  StatusInCourierWarehouse,
	97:  StatusInCourierWarehouse,
	99:  StatusInCourierWarehouse,
	100: StatusInCourierWarehouse,
	101: StatusInCourierWarehouse,
	103: StatusInCourierWarehouse,
	106: StatusInCourierWarehouse,
	105: StatusInCourierWarehouse,
	108: StatusInCourierWarehouse,
	110: StatusInCourierWarehouse,
	111: StatusInCourierWarehouse,
	129: StatusInCourierWarehouse,
	141: StatusInCourierWarehouse,
	149: StatusInCourierWarehouse,
	233: StatusInCourierWarehouse,
	249: StatusInCourierWarehouse,
	251: StatusInCourierWarehouse,
	255: StatusInCourierWarehouse,

	5:   StatusInTransit,
	24:  StatusInTransit,
	25:  StatusInTransit,
	27:  StatusInTransit,
	32:  StatusInTransit,
	34:  StatusInTransit,
	36:  StatusInTransit,
	146: StatusInTransit,
	200: StatusInTransit,
	203: StatusInTransit,
	204: StatusInTransit,
	206: StatusInTransit,
	208: StatusInTransit,
	212: StatusInTransit,
	213: StatusInTransit,
	216: StatusInTransit,

	120: StatusInTransit,
	114: StatusInTransit,
	121: StatusInTransit,
	122: StatusInTransit,
	124: StatusInTransit,
	128: StatusInTransit,
	127: StatusInTransit,
	132: StatusInTransit,
	133: StatusInTransit,
	140: StatusInTransit,
	145: StatusInTransit,
	199: StatusInTransit,

	21: StatusDelivered,
	60: StatusDelivered,
	64: StatusDelivered,

	153: StatusUndelivered,
	158: StatusUndelivered,
	159: StatusUndelivered,
	161: StatusUndelivered,
	180: StatusUndelivered,
	181: StatusUndelivered,
	182: StatusUndelivered,
	185: StatusUndelivered,
	186: StatusUndelivered,
	196: StatusUndelivered,
	197: StatusUndelivered,
	198: StatusUndelivered,
	96:  StatusUndelivered,
	54:  StatusUndelivered,
	89:  StatusUndelivered,
}

type UrgentCargusAwbInfoAdapter struct {
	BaseAdapter
}

/**
 * Transform
 * @param infoData interface{}|*mapper.CargusAwbInformationResponse
 * @param awbStatusDetails AwbStatusDetails pointer
 * @return AwbStatusDetails pointer
 * @desc - transforms the awb info from urgent cargus into a common structure
 */
func (adapter *UrgentCargusAwbInfoAdapter) Transform(infoData interface{}, awbStatusDetails *AwbStatusDetails) *AwbStatusDetails {
	info := infoData.(*mapper.CargusAwbInformationResponse)
	if info == nil {
		return awbStatusDetails
	}

	for _, history := range info.Event {
		// ignore some events
		if utilLib.InArray(ignoredEventIds, history.EventID) {
			continue
		}

		if urgentCargusEventMap[history.EventID] > 0 {
			awbStatusDetails.Status = urgentCargusEventMap[history.EventID]
		}

		awbStatusDetails.History = append(awbStatusDetails.History, StatusHistory{
			Status: history.Description,
			Date:   util.StringToDateTimeFormat(history.Date),
		})
	}

	return awbStatusDetails
}
