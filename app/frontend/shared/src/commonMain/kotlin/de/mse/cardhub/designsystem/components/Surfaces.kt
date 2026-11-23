package de.mse.cardhub.designsystem.components

import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.BoxScope
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.ColumnScope
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.material3.HorizontalDivider
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Brush
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.Dp
import androidx.compose.ui.unit.dp
import de.mse.cardhub.designsystem.icons.CardwiseIcon
import de.mse.cardhub.designsystem.icons.CardwiseIconImage
import de.mse.cardhub.designsystem.theme.CardwiseTheme

/** White surface with border. Base for panels like the PricePanel. */
@Composable
fun Panel(
    modifier: Modifier = Modifier,
    contentPadding: Dp = 15.dp,
    onClick: (() -> Unit)? = null,
    content: @Composable ColumnScope.() -> Unit,
) {
    val colors = CardwiseTheme.colors
    val shape = CardwiseTheme.shapes.card
    Column(
        modifier = modifier
            .fillMaxWidth()
            .clip(shape)
            .background(colors.surface, shape)
            .border(1.dp, colors.line, shape)
            .then(if (onClick != null) Modifier.clickable(onClick = onClick) else Modifier)
            .padding(contentPadding),
        content = content,
    )
}

/** Light blue info box. */
@Composable
fun InfoBanner(
    text: String,
    modifier: Modifier = Modifier,
) {
    Text(
        text = text,
        style = CardwiseTheme.typography.label.copy(lineHeight = CardwiseTheme.typography.label.fontSize * 1.6),
        color = CardwiseTheme.colors.blue,
        modifier = modifier
            .fillMaxWidth()
            .background(CardwiseTheme.colors.blueSoft, CardwiseTheme.shapes.lg)
            .padding(13.dp),
    )
}

/** Row with title, subtitle and a value on the right. */
@Composable
fun ValueRow(
    title: String,
    value: String,
    modifier: Modifier = Modifier,
    subtitle: String? = null,
) {
    val colors = CardwiseTheme.colors
    Row(modifier.fillMaxWidth(), horizontalArrangement = Arrangement.spacedBy(13.dp), verticalAlignment = Alignment.CenterVertically) {
        Column(Modifier.weight(1f)) {
            Text(title, style = CardwiseTheme.typography.bodyStrong, color = colors.text)
            if (subtitle != null) {
                Text(subtitle, style = CardwiseTheme.typography.label, color = colors.muted, modifier = Modifier.padding(top = 5.dp))
            }
        }
        Text(value, style = CardwiseTheme.typography.section, color = colors.text)
    }
}

/** Tile showing a single stat. */
@Composable
fun StatCard(
    value: String,
    label: String,
    modifier: Modifier = Modifier,
) {
    val colors = CardwiseTheme.colors
    val shape = CardwiseTheme.shapes.xl
    Column(
        modifier = modifier
            .background(colors.surface, shape)
            .border(1.dp, colors.line, shape)
            .padding(horizontal = 9.dp, vertical = 15.dp),
        horizontalAlignment = Alignment.CenterHorizontally,
    ) {
        Text(value, style = CardwiseTheme.typography.titleSmall, color = colors.text, textAlign = TextAlign.Center)
        Text(label, style = CardwiseTheme.typography.caption, color = colors.muted, textAlign = TextAlign.Center, modifier = Modifier.padding(top = 5.dp))
    }
}

/** Several stats side by side. */
@Composable
fun StatsRow(
    stats: List<Pair<String, String>>,
    modifier: Modifier = Modifier,
) {
    Row(modifier.fillMaxWidth(), horizontalArrangement = Arrangement.spacedBy(8.dp)) {
        stats.forEach { (value, label) -> StatCard(value, label, Modifier.weight(1f)) }
    }
}

/** Round avatar with initials, gradient and optional online dot. */
@Composable
fun Avatar(
    initials: String,
    modifier: Modifier = Modifier,
    size: Dp = 38.dp,
    gradient: List<Color> = AvatarGradients.Peach,
    online: Boolean = false,
) {
    val colors = CardwiseTheme.colors
    Box(modifier.size(size)) {
        Box(
            modifier = Modifier
                .size(size)
                .clip(CircleShape)
                .background(Brush.linearGradient(gradient))
                .border(2.dp, Color.White.copy(alpha = 0.75f), CircleShape),
            contentAlignment = Alignment.Center,
        ) {
            Text(
                initials,
                style = if (size >= 60.dp) CardwiseTheme.typography.title else CardwiseTheme.typography.labelStrong,
                color = colors.surface,
            )
        }
        if (online) {
            Box(
                Modifier
                    .align(Alignment.BottomEnd)
                    .size(size * 0.29f)
                    .background(colors.online, CircleShape)
                    .border(2.dp, colors.surface, CircleShape),
            )
        }
    }
}

object AvatarGradients {
    val Peach = listOf(Color(0xFFEEAE8A), Color(0xFFB96391))
    val Blue = listOf(Color(0xFF8DB4FF), Color(0xFF266DFF))
    val Mint = listOf(Color(0xFF8EE3C2), Color(0xFF2C9C8D))
    val Lilac = listOf(Color(0xFFC9B2FF), Color(0xFF7A5CE0))
}

/** Row for a person: avatar, name, info and optional trailing content. */
@Composable
fun PersonRow(
    name: String,
    subtitle: String,
    avatar: @Composable () -> Unit,
    modifier: Modifier = Modifier,
    trailing: (@Composable () -> Unit)? = null,
) {
    val colors = CardwiseTheme.colors
    Row(modifier.fillMaxWidth(), horizontalArrangement = Arrangement.spacedBy(10.dp), verticalAlignment = Alignment.CenterVertically) {
        avatar()
        Column(Modifier.weight(1f)) {
            Text(name, style = CardwiseTheme.typography.labelStrong, color = colors.text)
            Text(subtitle, style = CardwiseTheme.typography.caption, color = colors.muted, modifier = Modifier.padding(top = 4.dp))
        }
        trailing?.invoke()
    }
}

/** Grouped settings list with dividers between rows. */
@Composable
fun SettingsGroup(
    modifier: Modifier = Modifier,
    content: @Composable SettingsGroupScope.() -> Unit,
) {
    val colors = CardwiseTheme.colors
    val shape = CardwiseTheme.shapes.xl
    Column(
        modifier
            .fillMaxWidth()
            .clip(shape)
            .background(colors.surface, shape)
            .border(1.dp, colors.line, shape),
    ) {
        SettingsGroupScope(this).content()
    }
}

class SettingsGroupScope internal constructor(private val column: ColumnScope) {
    private var first = true

    /** A single row. Use [trailing] to show e.g. a [CardwiseSwitch] instead of the arrow. */
    @Composable
    fun SettingsRow(
        title: String,
        icon: CardwiseIcon,
        onClick: (() -> Unit)? = null,
        trailing: (@Composable () -> Unit)? = null,
    ) {
        val colors = CardwiseTheme.colors
        if (!first) HorizontalDivider(color = colors.line, thickness = 1.dp)
        first = false
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .then(if (onClick != null) Modifier.clickable(onClick = onClick) else Modifier)
                .padding(horizontal = 12.dp, vertical = 15.dp),
            horizontalArrangement = Arrangement.spacedBy(10.dp),
            verticalAlignment = Alignment.CenterVertically,
        ) {
            Box(Modifier.size(28.dp), contentAlignment = Alignment.CenterStart) {
                CardwiseIconImage(icon, tint = colors.blue, size = 18.dp)
            }
            Text(title, style = CardwiseTheme.typography.labelStrong, color = colors.text, modifier = Modifier.weight(1f))
            if (trailing != null) trailing() else CardwiseIconImage(CardwiseIcon.ChevronRight, tint = colors.muted, size = 16.dp)
        }
    }
}

/** Empty state, e.g. when a list has no items. */
@Composable
fun EmptyState(
    icon: CardwiseIcon,
    title: String,
    message: String,
    modifier: Modifier = Modifier,
    action: (@Composable () -> Unit)? = null,
) {
    val colors = CardwiseTheme.colors
    Column(
        modifier = modifier.fillMaxWidth().padding(vertical = 36.dp, horizontal = 24.dp),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.spacedBy(8.dp),
    ) {
        Box(Modifier.size(48.dp).background(colors.blueSoft, CardwiseTheme.shapes.lg), contentAlignment = Alignment.Center) {
            CardwiseIconImage(icon, tint = colors.blue, size = 22.dp)
        }
        Text(title, style = CardwiseTheme.typography.section, color = colors.text, textAlign = TextAlign.Center, modifier = Modifier.padding(top = 6.dp))
        Text(message, style = CardwiseTheme.typography.label, color = colors.muted, textAlign = TextAlign.Center)
        if (action != null) Box(Modifier.padding(top = 8.dp)) { action() }
    }
}

/** Action bar fixed to the bottom of the screen with a soft gradient on top. */
@Composable
fun BoxScope.StickyActionBar(
    modifier: Modifier = Modifier,
    content: @Composable () -> Unit,
) {
    val screen = CardwiseTheme.colors.screen
    Box(
        modifier = modifier
            .align(Alignment.BottomCenter)
            .fillMaxWidth()
            .background(Brush.verticalGradient(0f to screen.copy(alpha = 0f), 0.17f to screen))
            .padding(top = 26.dp, bottom = 22.dp),
    ) {
        content()
    }
}
