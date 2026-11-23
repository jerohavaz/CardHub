package de.mse.cardhub.designsystem.components

import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.defaultMinSize
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.alpha
import androidx.compose.ui.draw.clip
import androidx.compose.ui.draw.shadow
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.Shape
import androidx.compose.ui.semantics.Role
import androidx.compose.ui.unit.Dp
import androidx.compose.ui.unit.dp
import de.mse.cardhub.designsystem.icons.CardwiseIcon
import de.mse.cardhub.designsystem.icons.CardwiseIconImage
import de.mse.cardhub.designsystem.theme.CardwiseTheme

enum class ButtonVariant {
    /** Filled blue – main action */
    Primary,

    /** Light blue – secondary action */
    Soft,

    /** Gray – cancel */
    Neutral,

    /** Red – destructive action (remove, log out) */
    Danger,
}

/** Default button (44 dp high). */
@Composable
fun CardwiseButton(
    text: String,
    onClick: () -> Unit,
    modifier: Modifier = Modifier,
    variant: ButtonVariant = ButtonVariant.Primary,
    icon: CardwiseIcon? = null,
    enabled: Boolean = true,
) {
    val colors = CardwiseTheme.colors
    val (container, content) = when (variant) {
        ButtonVariant.Primary -> colors.blue to colors.surface
        ButtonVariant.Soft -> colors.blueSoft to colors.blue
        ButtonVariant.Neutral -> colors.surface2 to colors.text
        ButtonVariant.Danger -> colors.dangerSoft to colors.danger
    }
    val shape = CardwiseTheme.shapes.lg
    Row(
        modifier = modifier
            .defaultMinSize(minHeight = 44.dp)
            .alpha(if (enabled) 1f else 0.5f)
            .clip(shape)
            .background(container, shape)
            .clickable(enabled = enabled, role = Role.Button, onClick = onClick)
            .padding(horizontal = 15.dp, vertical = 10.dp),
        horizontalArrangement = Arrangement.spacedBy(8.dp, Alignment.CenterHorizontally),
        verticalAlignment = Alignment.CenterVertically,
    ) {
        if (icon != null) CardwiseIconImage(icon, tint = content, size = 16.dp)
        Text(text, style = CardwiseTheme.typography.button, color = content)
    }
}

/** Small dark pill button, used in the hero banner. */
@Composable
fun PillButton(
    text: String,
    onClick: () -> Unit,
    modifier: Modifier = Modifier,
) {
    val colors = CardwiseTheme.colors
    Box(
        modifier = modifier
            .height(31.dp)
            .clip(CardwiseTheme.shapes.pill)
            .background(colors.text)
            .clickable(role = Role.Button, onClick = onClick)
            .padding(horizontal = 14.dp),
        contentAlignment = Alignment.Center,
    ) {
        Text(text, style = CardwiseTheme.typography.captionStrong, color = colors.screen)
    }
}

enum class IconButtonVariant {
    /** White with border – back, favorite, close */
    Surface,

    /** Light blue – header actions, "+" on tiles */
    Soft,
}

/** Square or round icon button. */
@Composable
fun CardwiseIconButton(
    icon: CardwiseIcon,
    onClick: () -> Unit,
    contentDescription: String?,
    modifier: Modifier = Modifier,
    variant: IconButtonVariant = IconButtonVariant.Surface,
    shape: Shape = CardwiseTheme.shapes.lg,
    size: Dp = 39.dp,
    iconSize: Dp = 19.dp,
) {
    val colors = CardwiseTheme.colors
    val (container, tint, border) = when (variant) {
        IconButtonVariant.Surface -> Triple(colors.surface, colors.text, BorderStroke(1.dp, colors.line))
        IconButtonVariant.Soft -> Triple(colors.blueSoft, colors.blue, null)
    }
    Box(
        modifier = modifier
            .size(size)
            .shadow(elevation = 1.dp, shape = shape, ambientColor = ShadowColor, spotColor = ShadowColor)
            .clip(shape)
            .background(container, shape)
            .then(if (border != null) Modifier.border(border, shape) else Modifier)
            .clickable(role = Role.Button, onClick = onClick),
        contentAlignment = Alignment.Center,
    ) {
        CardwiseIconImage(icon, tint = tint, size = iconSize, contentDescription = contentDescription)
    }
}

/** Round back button for the [BackHeader]. */
@Composable
fun CircleIconButton(
    icon: CardwiseIcon,
    onClick: () -> Unit,
    contentDescription: String?,
    modifier: Modifier = Modifier,
) = CardwiseIconButton(
    icon = icon,
    onClick = onClick,
    contentDescription = contentDescription,
    modifier = modifier,
    shape = CircleShape,
    size = 37.dp,
)

/** Small "+" button on card tiles. */
@Composable
fun AddButton(
    onClick: () -> Unit,
    contentDescription: String?,
    modifier: Modifier = Modifier,
) = CardwiseIconButton(
    icon = CardwiseIcon.Plus,
    onClick = onClick,
    contentDescription = contentDescription,
    modifier = modifier,
    variant = IconButtonVariant.Soft,
    shape = RoundedCornerShape(9.dp),
    size = 30.dp,
    iconSize = 16.dp,
)

internal val ShadowColor = Color(0x14233655)
