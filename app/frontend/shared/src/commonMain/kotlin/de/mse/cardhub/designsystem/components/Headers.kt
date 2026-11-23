package de.mse.cardhub.designsystem.components

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.RowScope
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import de.mse.cardhub.designsystem.icons.CardwiseIcon
import de.mse.cardhub.designsystem.theme.CardwiseTheme

/** Header of the main tabs: overline, title, subtitle and an optional action on the right. */
@Composable
fun ScreenHeader(
    eyebrow: String,
    title: String,
    modifier: Modifier = Modifier,
    subtitle: String? = null,
    action: (@Composable () -> Unit)? = null,
) {
    val colors = CardwiseTheme.colors
    val type = CardwiseTheme.typography
    Row(
        modifier = modifier.fillMaxWidth().padding(top = 15.dp),
        verticalAlignment = Alignment.Top,
    ) {
        Column(Modifier.weight(1f)) {
            Text(eyebrow.uppercase(), style = type.eyebrow, color = colors.muted)
            Text(title, style = type.display, color = colors.text, modifier = Modifier.padding(top = 7.dp))
            if (subtitle != null) {
                Text(subtitle, style = type.bodySmall, color = colors.muted, modifier = Modifier.padding(top = 8.dp))
            }
        }
        action?.invoke()
    }
}

/** Header for subpages: round back button, title and optional actions on the right. */
@Composable
fun BackHeader(
    title: String,
    onBack: () -> Unit,
    modifier: Modifier = Modifier,
    actions: @Composable RowScope.() -> Unit = {},
) {
    Row(
        modifier = modifier.fillMaxWidth().padding(top = 8.dp),
        horizontalArrangement = Arrangement.spacedBy(11.dp),
        verticalAlignment = Alignment.CenterVertically,
    ) {
        CircleIconButton(CardwiseIcon.ChevronLeft, onClick = onBack, contentDescription = "Zurück")
        Text(
            text = title,
            style = CardwiseTheme.typography.titleSmall,
            color = CardwiseTheme.colors.text,
            modifier = Modifier.weight(1f),
        )
        actions()
    }
}

/** Section heading with an optional link on the right. */
@Composable
fun SectionHeader(
    title: String,
    modifier: Modifier = Modifier,
    actionText: String? = null,
    onAction: (() -> Unit)? = null,
) {
    Row(
        modifier = modifier.fillMaxWidth().padding(top = 20.dp),
        verticalAlignment = Alignment.CenterVertically,
    ) {
        Text(title, style = CardwiseTheme.typography.section, color = CardwiseTheme.colors.text, modifier = Modifier.weight(1f))
        if (actionText != null) {
            Text(
                text = actionText,
                style = CardwiseTheme.typography.caption,
                color = CardwiseTheme.colors.muted,
                modifier = if (onAction != null) Modifier.clickable(onClick = onAction) else Modifier,
            )
        }
    }
}
