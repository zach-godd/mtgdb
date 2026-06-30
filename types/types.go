package types

// Card represents a single MTG card object from Scryfall.
type Card struct {
	ID              string        `json:"id"               bson:"_id"`
	Object          string        `json:"object"           bson:"object"`
	OracleID        string        `json:"oracle_id"        bson:"oracle_id,omitempty"`
	MultiverseIDs   []int         `json:"multiverse_ids"   bson:"multiverse_ids"`
	MtgoID          *int          `json:"mtgo_id"          bson:"mtgo_id,omitempty"`
	MtgoFoilID      *int          `json:"mtgo_foil_id"     bson:"mtgo_foil_id,omitempty"`
	ArenaID         *int          `json:"arena_id"         bson:"arena_id,omitempty"`
	TcgplayerID     *int          `json:"tcgplayer_id"     bson:"tcgplayer_id,omitempty"`
	CardmarketID    *int          `json:"cardmarket_id"    bson:"cardmarket_id,omitempty"`
	Name            string        `json:"name"             bson:"name"`
	PrintedName     *string       `json:"printed_name"     bson:"printed_name,omitempty"`
	Lang            string        `json:"lang"             bson:"lang"`
	ReleasedAt      string        `json:"released_at"      bson:"released_at"`
	URI             string        `json:"uri"              bson:"uri"`
	ScryfallURI     string        `json:"scryfall_uri"     bson:"scryfall_uri"`
	Layout          string        `json:"layout"           bson:"layout"`
	HighresImage    bool          `json:"highres_image"    bson:"highres_image"`
	ImageStatus     string        `json:"image_status"     bson:"image_status"`
	ImageURIs       *ImageURIs    `json:"image_uris"       bson:"image_uris,omitempty"`
	ManaCost        string        `json:"mana_cost"        bson:"mana_cost,omitempty"`
	CMC             float64       `json:"cmc"              bson:"cmc"`
	TypeLine        string        `json:"type_line"        bson:"type_line"`
	PrintedTypeLine *string       `json:"printed_type_line" bson:"printed_type_line,omitempty"`
	OracleText      string        `json:"oracle_text"      bson:"oracle_text,omitempty"`
	PrintedText     *string       `json:"printed_text"     bson:"printed_text,omitempty"`
	FlavorText      *string       `json:"flavor_text"      bson:"flavor_text,omitempty"`
	FlavorName      *string       `json:"flavor_name"      bson:"flavor_name,omitempty"`
	Power           *string       `json:"power"            bson:"power,omitempty"`
	Toughness       *string       `json:"toughness"        bson:"toughness,omitempty"`
	Loyalty         *string       `json:"loyalty"          bson:"loyalty,omitempty"`
	Defense         *string       `json:"defense"          bson:"defense,omitempty"`
	HandModifier    *string       `json:"hand_modifier"    bson:"hand_modifier,omitempty"`
	LifeModifier    *string       `json:"life_modifier"    bson:"life_modifier,omitempty"`
	Colors          []string      `json:"colors"           bson:"colors"`
	ColorIdentity   []string      `json:"color_identity"   bson:"color_identity"`
	ColorIndicator  []string      `json:"color_indicator"  bson:"color_indicator,omitempty"`
	Keywords        []string      `json:"keywords"         bson:"keywords"`
	ProducedMana    []string      `json:"produced_mana"    bson:"produced_mana,omitempty"`
	AllParts        []RelatedCard `json:"all_parts"     bson:"all_parts,omitempty"`
	CardFaces       []CardFace    `json:"card_faces"       bson:"card_faces,omitempty"`
	Legalities      Legalities    `json:"legalities"       bson:"legalities"`
	Games           []string      `json:"games"            bson:"games"`
	Reserved        bool          `json:"reserved"         bson:"reserved"`
	GameChanger     bool          `json:"game_changer"     bson:"game_changer"`
	Foil            bool          `json:"foil"             bson:"foil"`
	Nonfoil         bool          `json:"nonfoil"          bson:"nonfoil"`
	Finishes        []string      `json:"finishes"         bson:"finishes"`
	Oversized       bool          `json:"oversized"        bson:"oversized"`
	Promo           bool          `json:"promo"            bson:"promo"`
	Reprint         bool          `json:"reprint"          bson:"reprint"`
	Variation       bool          `json:"variation"        bson:"variation"`
	VariationOf     *string       `json:"variation_of"     bson:"variation_of,omitempty"`
	SetID           string        `json:"set_id"           bson:"set_id"`
	Set             string        `json:"set"              bson:"set"`
	SetName         string        `json:"set_name"         bson:"set_name"`
	SetType         string        `json:"set_type"         bson:"set_type"`
	SetURI          string        `json:"set_uri"          bson:"set_uri"`
	SetSearchURI    string        `json:"set_search_uri"   bson:"set_search_uri"`
	ScryfallSetURI  string        `json:"scryfall_set_uri" bson:"scryfall_set_uri"`
	RulingsURI      string        `json:"rulings_uri"      bson:"rulings_uri"`
	PrintsSearchURI string        `json:"prints_search_uri" bson:"prints_search_uri"`
	CollectorNumber string        `json:"collector_number" bson:"collector_number"`
	Digital         bool          `json:"digital"          bson:"digital"`
	Rarity          string        `json:"rarity"           bson:"rarity"`
	CardBackID      *string       `json:"card_back_id"     bson:"card_back_id,omitempty"`
	Artist          *string       `json:"artist"           bson:"artist,omitempty"`
	ArtistIDs       []string      `json:"artist_ids"       bson:"artist_ids,omitempty"`
	IllustrationID  *string       `json:"illustration_id"  bson:"illustration_id,omitempty"`
	BorderColor     string        `json:"border_color"     bson:"border_color"`
	Frame           string        `json:"frame"            bson:"frame"`
	FrameEffects    []string      `json:"frame_effects"    bson:"frame_effects,omitempty"`
	SecurityStamp   *string       `json:"security_stamp"   bson:"security_stamp,omitempty"`
	FullArt         bool          `json:"full_art"         bson:"full_art"`
	Textless        bool          `json:"textless"         bson:"textless"`
	Booster         bool          `json:"booster"          bson:"booster"`
	StorySpotlight  bool          `json:"story_spotlight"  bson:"story_spotlight"`
	EdhrecRank      *int          `json:"edhrec_rank"      bson:"edhrec_rank,omitempty"`
	PennyRank       *int          `json:"penny_rank"       bson:"penny_rank,omitempty"`
	PromoTypes      []string      `json:"promo_types"      bson:"promo_types,omitempty"`
	Watermark       *string       `json:"watermark"        bson:"watermark,omitempty"`
	Prices          Prices        `json:"prices"           bson:"prices"`
	RelatedURIs     RelatedURIs   `json:"related_uris"   bson:"related_uris"`
	PurchaseURIs    *PurchaseURIs `json:"purchase_uris" bson:"purchase_uris,omitempty"`
}

type ImageURIs struct {
	Small      string `json:"small"        bson:"small"`
	Normal     string `json:"normal"       bson:"normal"`
	Large      string `json:"large"        bson:"large"`
	PNG        string `json:"png"          bson:"png"`
	ArtCrop    string `json:"art_crop"     bson:"art_crop"`
	BorderCrop string `json:"border_crop"  bson:"border_crop"`
}

// CardFace represents one face of a multi-faced card (transform, adventure, split, etc.).
type CardFace struct {
	Object          string     `json:"object"            bson:"object"`
	Name            string     `json:"name"              bson:"name"`
	PrintedName     *string    `json:"printed_name"      bson:"printed_name,omitempty"`
	ManaCost        string     `json:"mana_cost"         bson:"mana_cost"`
	TypeLine        string     `json:"type_line"         bson:"type_line,omitempty"`
	PrintedTypeLine *string    `json:"printed_type_line" bson:"printed_type_line,omitempty"`
	OracleText      string     `json:"oracle_text"       bson:"oracle_text,omitempty"`
	PrintedText     *string    `json:"printed_text"      bson:"printed_text,omitempty"`
	FlavorText      *string    `json:"flavor_text"       bson:"flavor_text,omitempty"`
	Power           *string    `json:"power"             bson:"power,omitempty"`
	Toughness       *string    `json:"toughness"         bson:"toughness,omitempty"`
	Loyalty         *string    `json:"loyalty"           bson:"loyalty,omitempty"`
	Defense         *string    `json:"defense"           bson:"defense,omitempty"`
	Colors          []string   `json:"colors"            bson:"colors,omitempty"`
	ColorIndicator  []string   `json:"color_indicator"   bson:"color_indicator,omitempty"`
	Watermark       *string    `json:"watermark"         bson:"watermark,omitempty"`
	Artist          *string    `json:"artist"            bson:"artist,omitempty"`
	ArtistID        *string    `json:"artist_id"         bson:"artist_id,omitempty"`
	IllustrationID  *string    `json:"illustration_id"   bson:"illustration_id,omitempty"`
	ImageURIs       *ImageURIs `json:"image_uris"        bson:"image_uris,omitempty"`
}

// RelatedCard is a reference to another card that shares a relationship with this one.
type RelatedCard struct {
	Object    string `json:"object"    bson:"object"`
	ID        string `json:"id"        bson:"id"`
	Component string `json:"component" bson:"component"`
	Name      string `json:"name"      bson:"name"`
	TypeLine  string `json:"type_line" bson:"type_line"`
	URI       string `json:"uri"       bson:"uri"`
}

type Legalities struct {
	Standard         string `json:"standard"          bson:"standard"`
	Future           string `json:"future"            bson:"future"`
	Historic         string `json:"historic"          bson:"historic"`
	Timeless         string `json:"timeless"          bson:"timeless"`
	Gladiator        string `json:"gladiator"         bson:"gladiator"`
	Pioneer          string `json:"pioneer"           bson:"pioneer"`
	Modern           string `json:"modern"            bson:"modern"`
	Legacy           string `json:"legacy"            bson:"legacy"`
	Pauper           string `json:"pauper"            bson:"pauper"`
	Vintage          string `json:"vintage"           bson:"vintage"`
	Penny            string `json:"penny"             bson:"penny"`
	Commander        string `json:"commander"         bson:"commander"`
	Oathbreaker      string `json:"oathbreaker"       bson:"oathbreaker"`
	StandardBrawl    string `json:"standardbrawl"     bson:"standardbrawl"`
	Brawl            string `json:"brawl"             bson:"brawl"`
	CompetitiveBrawl string `json:"competitivebrawl"  bson:"competitivebrawl"`
	Alchemy          string `json:"alchemy"           bson:"alchemy"`
	PauperCommander  string `json:"paupercommander"   bson:"paupercommander"`
	Duel             string `json:"duel"              bson:"duel"`
	OldSchool        string `json:"oldschool"         bson:"oldschool"`
	Premodern        string `json:"premodern"         bson:"premodern"`
	PreDH            string `json:"predh"             bson:"predh"`
	TLR              string `json:"tlr"               bson:"tlr"`
}

// Prices are nullable strings because Scryfall returns JSON null when no price exists.
type Prices struct {
	USD       *string `json:"usd"        bson:"usd"`
	USDFoil   *string `json:"usd_foil"   bson:"usd_foil"`
	USDEtched *string `json:"usd_etched" bson:"usd_etched"`
	EUR       *string `json:"eur"        bson:"eur"`
	EURFoil   *string `json:"eur_foil"   bson:"eur_foil"`
	Tix       *string `json:"tix"        bson:"tix"`
}

type RelatedURIs struct {
	Gatherer                  string `json:"gatherer"                    bson:"gatherer,omitempty"`
	TcgplayerInfiniteArticles string `json:"tcgplayer_infinite_articles" bson:"tcgplayer_infinite_articles,omitempty"`
	TcgplayerInfiniteDecks    string `json:"tcgplayer_infinite_decks"    bson:"tcgplayer_infinite_decks,omitempty"`
	Edhrec                    string `json:"edhrec"                      bson:"edhrec,omitempty"`
	Mtgtop8                   string `json:"mtgtop8"                     bson:"mtgtop8,omitempty"`
}

type PurchaseURIs struct {
	Tcgplayer   string `json:"tcgplayer"   bson:"tcgplayer,omitempty"`
	Cardmarket  string `json:"cardmarket"  bson:"cardmarket,omitempty"`
	Cardhoarder string `json:"cardhoarder" bson:"cardhoarder,omitempty"`
}
