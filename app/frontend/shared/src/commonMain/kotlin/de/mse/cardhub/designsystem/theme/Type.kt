package de.mse.cardhub.designsystem.theme

import androidx.compose.runtime.Immutable
import androidx.compose.ui.text.TextStyle
import androidx.compose.ui.text.font.FontFamily
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.em
import androidx.compose.ui.unit.sp

/**
 * Typography from the Figma template. It uses SF Pro – we use the
 * system font (iOS uses SF Pro and Android uses Roboto) to avoid licensing issue
 */
@Immutable
data class CardHubTypography(
    /** Screen title, e.g. "Pokédex" */
    val display: TextStyle,
    /** Card name in the details view, e.g. "Auralith" */
    val headline: TextStyle,
    /** Hero headline "Entdecke das Lumen-Set." */
    val hero: TextStyle,
    /** Large price in the PricePanel */
    val price: TextStyle,
    /** Dialog/sheet title */
    val title: TextStyle,
    /** Title in the BackHeader, stat values */
    val titleSmall: TextStyle,
    /** Section headings */
    val section: TextStyle,
    val bodyStrong: TextStyle,
    val body: TextStyle,
    val bodySmall: TextStyle,
    val label: TextStyle,
    val labelStrong: TextStyle,
    val caption: TextStyle,
    val captionStrong: TextStyle,
    /** Overline "CARDHUB / POKÉDEX" */
    val eyebrow: TextStyle,
    val button: TextStyle,
    val badge: TextStyle,
    val badgeSmall: TextStyle,
)

private val Base = TextStyle(fontFamily = FontFamily.Default)

val DefaultCardHubTypography = CardHubTypography(
    display = Base.copy(fontSize = 30.sp, lineHeight = 33.sp, fontWeight = FontWeight.Bold, letterSpacing = (-0.045).em),
    headline = Base.copy(fontSize = 26.sp, lineHeight = 39.sp, fontWeight = FontWeight.Bold, letterSpacing = (-0.035).em),
    hero = Base.copy(fontSize = 23.sp, lineHeight = 25.sp, fontWeight = FontWeight.Bold, letterSpacing = (-0.043).em),
    price = Base.copy(fontSize = 27.sp, lineHeight = 40.sp, fontWeight = FontWeight.Bold, letterSpacing = (-0.03).em),
    title = Base.copy(fontSize = 20.sp, lineHeight = 30.sp, fontWeight = FontWeight.Bold, letterSpacing = (-0.025).em),
    titleSmall = Base.copy(fontSize = 19.sp, lineHeight = 28.5.sp, fontWeight = FontWeight.Bold, letterSpacing = (-0.026).em),
    section = Base.copy(fontSize = 16.sp, lineHeight = 24.sp, fontWeight = FontWeight.Bold),
    bodyStrong = Base.copy(fontSize = 13.sp, lineHeight = 16.sp, fontWeight = FontWeight.Bold),
    body = Base.copy(fontSize = 13.sp, lineHeight = 16.sp),
    bodySmall = Base.copy(fontSize = 12.sp, lineHeight = 18.sp),
    label = Base.copy(fontSize = 11.sp, lineHeight = 16.5.sp),
    labelStrong = Base.copy(fontSize = 12.sp, lineHeight = 18.sp, fontWeight = FontWeight.Bold),
    caption = Base.copy(fontSize = 10.sp, lineHeight = 15.sp),
    captionStrong = Base.copy(fontSize = 10.sp, lineHeight = 15.sp, fontWeight = FontWeight.ExtraBold),
    eyebrow = Base.copy(fontSize = 10.sp, lineHeight = 15.sp, fontWeight = FontWeight.ExtraBold, letterSpacing = 0.15.em),
    button = Base.copy(fontSize = 12.sp, lineHeight = 18.sp, fontWeight = FontWeight.ExtraBold),
    badge = Base.copy(fontSize = 9.sp, lineHeight = 13.5.sp, fontWeight = FontWeight.ExtraBold),
    badgeSmall = Base.copy(fontSize = 8.sp, lineHeight = 12.sp, fontWeight = FontWeight.ExtraBold),
)
