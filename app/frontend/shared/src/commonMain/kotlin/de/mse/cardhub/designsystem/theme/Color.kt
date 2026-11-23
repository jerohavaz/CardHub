package de.mse.cardhub.designsystem.theme

import androidx.compose.runtime.Immutable
import androidx.compose.ui.graphics.Color

/**
 * Color tokens from the Figma template "CARDWISE · Glacier".
 * The names match the CSS variables in the design (--screen, --surface, --text …).
 */
@Immutable
data class CardwiseColors(
    val screen: Color,
    val surface: Color,
    val surface2: Color,
    val line: Color,
    val text: Color,
    val muted: Color,
    val placeholder: Color,
    val navInactive: Color,
    val blue: Color,
    val blueSoft: Color,
    val heroText: Color,
    val success: Color,
    val online: Color,
    val danger: Color,
    val dangerSoft: Color,
    val dangerIcon: Color,
    val warning: Color,
    val warningSoft: Color,
    val gold: Color,
    val goldSoft: Color,
    val neutralSoft: Color,
)

val LightCardwiseColors = CardwiseColors(
    screen = Color(0xFFF8FAFF),
    surface = Color(0xFFFFFFFF),
    surface2 = Color(0xFFEEF3FB),
    line = Color(0xFFE7ECF4),
    text = Color(0xFF111E38),
    muted = Color(0xFF8290A9),
    placeholder = Color(0xFF9FAABD),
    navInactive = Color(0xFF91A0BC),
    blue = Color(0xFF266DFF),
    blueSoft = Color(0xFFE8F0FF),
    heroText = Color(0xFF14264C),
    success = Color(0xFF20A474),
    online = Color(0xFF2CCC8D),
    danger = Color(0xFFE15262),
    dangerSoft = Color(0x17E65468),
    dangerIcon = Color(0xFFEC6B76),
    warning = Color(0xFFCA8A28),
    warningSoft = Color(0xFFFFF3DD),
    gold = Color(0xFF966900),
    goldSoft = Color(0xFFFFF4CC),
    neutralSoft = Color(0xFFF1F5FC),
)
