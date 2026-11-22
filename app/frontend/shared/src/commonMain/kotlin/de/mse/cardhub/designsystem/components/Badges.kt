package de.mse.cardhub.designsystem.components

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.unit.dp
import de.mse.cardhub.designsystem.theme.CardHubColors
import de.mse.cardhub.designsystem.theme.CardHubTheme

enum class BadgeTone { Blue, Warning, Gold, Neutral, Success, Danger }

private fun CardHubColors.badgeColors(tone: BadgeTone): Pair<Color, Color> = when (tone) {
    BadgeTone.Blue -> blueSoft to blue
    BadgeTone.Warning -> warningSoft to warning
    BadgeTone.Gold -> goldSoft to gold
    BadgeTone.Neutral -> neutralSoft to muted
    BadgeTone.Success -> Color(0xFFE3F6EE) to success
    BadgeTone.Danger -> dangerSoft to danger
}

/** Pill-shaped badge, e.g. for rarity, prices or trade status. */
@Composable
fun CardHubBadge(
    text: String,
    modifier: Modifier = Modifier,
    tone: BadgeTone = BadgeTone.Blue,
) {
    val (container, content) = CardHubTheme.colors.badgeColors(tone)
    Text(
        text = text,
        style = CardHubTheme.typography.badge,
        color = content,
        modifier = modifier
            .background(container, CardHubTheme.shapes.pill)
            .padding(horizontal = 8.dp, vertical = 5.dp),
    )
}

/** Card rarity. Labels match the texts in the design. */
enum class Rarity(val label: String, val tone: BadgeTone) {
    Common("Häufig", BadgeTone.Neutral),
    Rare("Selten", BadgeTone.Gold),
    Ultra("Ultra", BadgeTone.Gold),
    Legendary("Legendär", BadgeTone.Gold),
}

/** Small rarity badge for card tiles. */
@Composable
fun RarityTag(
    rarity: Rarity,
    modifier: Modifier = Modifier,
) {
    val (container, content) = CardHubTheme.colors.badgeColors(rarity.tone)
    Text(
        text = rarity.label,
        style = CardHubTheme.typography.badgeSmall,
        color = content,
        modifier = modifier
            .background(container, CardHubTheme.shapes.xs)
            .padding(horizontal = 6.dp, vertical = 4.dp),
    )
}
