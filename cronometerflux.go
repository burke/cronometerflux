package cronometerflux

import (
    "fmt"
    "strings"
    "github.com/burke/gocronometer"
)

// nutrientInfo holds metadata about how to export a nutrient
type nutrientInfo struct {
    value float64
    unit  string
}

// getNutrientMap extracts all nutrients from a serving record
func getNutrientMap(s gocronometer.ServingRecord) map[string]nutrientInfo {
    return map[string]nutrientInfo{
        // Energy and Water
        "energy":   {s.EnergyKcal, "kcal"},
        "water":    {s.WaterG, "g"},
        "caffeine": {s.CaffeineMg, "mg"},

        // B Vitamins
        "vitamin_b1":  {s.B1Mg, "mg"},
        "vitamin_b2":  {s.B2Mg, "mg"},
        "vitamin_b3":  {s.B3Mg, "mg"},
        "vitamin_b5":  {s.B5Mg, "mg"},
        "vitamin_b6":  {s.B6Mg, "mg"},
        "vitamin_b12": {s.B12Mg, "µg"},
        "biotin":      {s.BiotinUg, "µg"},
        "choline":     {s.CholineMg, "mg"},
        "folate":      {s.FolateUg, "µg"},

        // Other Vitamins
        "vitamin_a": {s.VitaminAUg, "µg"},
        "vitamin_c": {s.VitaminCMg, "mg"},
        "vitamin_d": {s.VitaminDUI, "IU"},
        "vitamin_e": {s.VitaminEMg, "mg"},
        "vitamin_k": {s.VitaminKMg, "µg"},

        // Minerals
        "calcium":    {s.CalciumMg, "mg"},
        "chromium":   {s.ChromiumUg, "µg"},
        "copper":     {s.CopperMg, "mg"},
        "fluoride":   {s.FluorideUg, "µg"},
        "iodine":     {s.IodineUg, "µg"},
        "iron":       {s.IronMg, "mg"},
        "magnesium":  {s.MagnesiumMg, "mg"},
        "manganese":  {s.ManganeseMg, "mg"},
        "phosphorus": {s.PhosphorusMg, "mg"},
        "potassium":  {s.PotassiumMg, "mg"},
        "selenium":   {s.SeleniumUg, "µg"},
        "sodium":     {s.SodiumMg, "mg"},
        "zinc":       {s.ZincMg, "mg"},

        // Carbohydrates
        "carbs":     {s.CarbsG, "g"},
        "fiber":     {s.FiberG, "g"},
        "fructose":  {s.FructoseG, "g"},
        "galactose": {s.GalactoseG, "g"},
        "glucose":   {s.GlucoseG, "g"},
        "lactose":   {s.LactoseG, "g"},
        "maltose":   {s.MaltoseG, "g"},
        "starch":    {s.StarchG, "g"},
        "sucrose":   {s.SucroseG, "g"},
        "sugars":    {s.SugarsG, "g"},
        "net_carbs": {s.NetCarbsG, "g"},

        // Fats
        "fat":             {s.FatG, "g"},
        "cholesterol":     {s.CholesterolMg, "mg"},
        "monounsaturated": {s.MonounsaturatedG, "g"},
        "polyunsaturated": {s.PolyunsaturatedG, "g"},
        "saturated":       {s.SaturatedG, "g"},
        "trans_fat":       {s.TransFatG, "g"},
        "omega3":          {s.Omega3G, "g"},
        "omega6":          {s.Omega6G, "g"},

        // Amino Acids
        "cystine":       {s.CystineG, "g"},
        "histidine":     {s.HistidineG, "g"},
        "isoleucine":    {s.IsoleucineG, "g"},
        "leucine":       {s.LeucineG, "g"},
        "lysine":        {s.LysineG, "g"},
        "methionine":    {s.MethionineG, "g"},
        "phenylalanine": {s.PhenylalanineG, "g"},
        "threonine":     {s.ThreonineG, "g"},
        "tryptophan":    {s.TryptophanG, "g"},
        "tyrosine":      {s.TyrosineG, "g"},
        "valine":        {s.ValineG, "g"},
        "protein":       {s.ProteinG, "g"},

        // Alcohol
        "alcohol":       {s.AlcoholG, "g"},
    }
}


// escapeTag handles escaping special characters in InfluxDB line protocol tag values
func escapeTag(s string) string {
    if s == "" {
        return "unknown"
    }
    // Ensure we escape the entire string properly
    s = strings.ReplaceAll(s, ",", "\\,")
    s = strings.ReplaceAll(s, "=", "\\=")
    s = strings.ReplaceAll(s, " ", "\\ ")
    return s
}

// escapeField handles escaping and quoting string field values
func escapeField(s string) string {
    if s == "" {
        return "\"\""
    }
    s = strings.ReplaceAll(s, "\"", "\\\"")
    return fmt.Sprintf("\"%s\"", s)
}

// FormatServing converts a single serving record to InfluxDB line protocol format
func FormatServing(serving gocronometer.ServingRecord) []string {
    var lines []string

    // Create base tags
    foodName := escapeTag(serving.FoodName)
    group := escapeTag(serving.Group)
    category := escapeTag(serving.Category)

    // Build the base tag string
    baseTags := fmt.Sprintf("food=%s,group=%s,category=%s",
        foodName, group, category)

    // Output quantity measurement
    lines = append(lines, fmt.Sprintf("nutrition_serving,%s,nutrient=quantity value=%.3f,units=%s %d",
        baseTags,
        serving.QuantityValue,
        escapeField(serving.QuantityUnits),
        serving.RecordedTime.UnixNano()))

    // Output all nutrients
    for nutrientName, info := range getNutrientMap(serving) {
        if info.value != 0 { // Skip zero values to reduce noise
            lines = append(lines, fmt.Sprintf("nutrition_nutrient,%s,nutrient=%s value=%.3f,units=%s %d",
                baseTags,
                escapeTag(nutrientName),
                info.value,
                escapeField(info.unit),
                serving.RecordedTime.UnixNano()))
        }
    }

    return lines
}

// FormatServings converts a slice of serving records to InfluxDB line protocol format
func FormatServings(servings gocronometer.ServingRecords) []string {
    var lines []string
    for _, serving := range servings {
        lines = append(lines, FormatServing(serving)...)
    }
    return lines
}
