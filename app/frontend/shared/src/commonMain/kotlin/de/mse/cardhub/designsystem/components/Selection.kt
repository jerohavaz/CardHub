package de.mse.cardhub.designsystem.components

import androidx.compose.animation.animateColorAsState
import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.PaddingValues
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.defaultMinSize
import androidx.compose.foundation.layout.fillMaxHeight
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.offset
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.lazy.LazyRow
import androidx.compose.foundation.lazy.itemsIndexed
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.draw.shadow
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.semantics.Role
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.Dp
import androidx.compose.ui.unit.dp
import de.mse.cardhub.designsystem.theme.CardwiseSpacing
import de.mse.cardhub.designsystem.theme.CardwiseTheme

/** Filter chip, e.g. for card types. */
@Composable
fun CardwiseChip(
    text: String,
    selected: Boolean,
    onClick: () -> Unit,
    modifier: Modifier = Modifier,
) {
    val colors = CardwiseTheme.colors
    val shape = CardwiseTheme.shapes.pill
    val container by animateColorAsState(if (selected) colors.blue else colors.surface)
    val border by animateColorAsState(if (selected) colors.blue else colors.line)
    Box(
        modifier = modifier
            .defaultMinSize(minHeight = 33.dp)
            .clip(shape)
            .background(container, shape)
            .border(1.dp, border, shape)
            .clickable(role = Role.Tab, onClick = onClick)
            .padding(horizontal = 12.dp, vertical = 8.dp),
        contentAlignment = Alignment.Center,
    ) {
        Text(
            text = text,
            style = CardwiseTheme.typography.label.copy(fontWeight = FontWeight.SemiBold),
            color = if (selected) colors.surface else colors.muted,
        )
    }
}

/** Horizontally scrollable row of chips that extends to the screen edge. */
@Composable
fun ChipRow(
    options: List<String>,
    selectedIndex: Int,
    onSelect: (Int) -> Unit,
    modifier: Modifier = Modifier,
    contentPadding: PaddingValues = PaddingValues(horizontal = CardwiseSpacing.screenHorizontal),
) {
    LazyRow(
        modifier = modifier,
        contentPadding = contentPadding,
        horizontalArrangement = Arrangement.spacedBy(7.dp),
    ) {
        itemsIndexed(options) { index, option ->
            CardwiseChip(option, selected = index == selectedIndex, onClick = { onSelect(index) })
        }
    }
}

/** Segmented control for switching between tabs. */
@Composable
fun SegmentedControl(
    options: List<String>,
    selectedIndex: Int,
    onSelect: (Int) -> Unit,
    modifier: Modifier = Modifier,
) {
    val colors = CardwiseTheme.colors
    Row(
        modifier = modifier
            .fillMaxWidth()
            .height(42.dp)
            .background(colors.surface2, RoundedCornerShape(13.dp))
            .padding(4.dp),
        horizontalArrangement = Arrangement.spacedBy(2.dp),
    ) {
        options.forEachIndexed { index, option ->
            val selected = index == selectedIndex
            val shape = CardwiseTheme.shapes.sm
            val container by animateColorAsState(if (selected) colors.surface else Color.Transparent)
            Box(
                modifier = Modifier
                    .weight(1f)
                    .fillMaxHeight()
                    .then(if (selected) Modifier.shadow(2.dp, shape, ambientColor = ShadowColor, spotColor = ShadowColor) else Modifier)
                    .clip(shape)
                    .background(container, shape)
                    .clickable(role = Role.Tab, onClick = { onSelect(index) }),
                contentAlignment = Alignment.Center,
            ) {
                Text(
                    text = option,
                    style = if (selected) CardwiseTheme.typography.captionStrong else CardwiseTheme.typography.caption,
                    color = if (selected) colors.text else colors.muted,
                )
            }
        }
    }
}

/** Quantity stepper (− 2 +). */
@Composable
fun Stepper(
    value: Int,
    onValueChange: (Int) -> Unit,
    modifier: Modifier = Modifier,
    min: Int = 1,
    max: Int = 99,
) {
    val colors = CardwiseTheme.colors
    Row(
        modifier = modifier
            .background(colors.surface2, CardwiseTheme.shapes.md)
            .padding(3.dp),
        horizontalArrangement = Arrangement.spacedBy(9.dp),
        verticalAlignment = Alignment.CenterVertically,
    ) {
        StepperButton("−", enabled = value > min, contentDescription = "Anzahl verringern") { onValueChange(value - 1) }
        Text(
            text = value.toString(),
            style = CardwiseTheme.typography.button,
            color = colors.text,
            modifier = Modifier.defaultMinSize(minWidth = 12.dp),
        )
        StepperButton("+", enabled = value < max, contentDescription = "Anzahl erhöhen") { onValueChange(value + 1) }
    }
}

@Composable
private fun StepperButton(symbol: String, enabled: Boolean, contentDescription: String, onClick: () -> Unit) {
    val colors = CardwiseTheme.colors
    Box(
        modifier = Modifier
            .size(28.dp)
            .clip(RoundedCornerShape(9.dp))
            .background(colors.surface)
            .clickable(enabled = enabled, role = Role.Button, onClickLabel = contentDescription, onClick = onClick),
        contentAlignment = Alignment.Center,
    ) {
        Text(
            symbol,
            style = CardwiseTheme.typography.section,
            color = if (enabled) colors.blue else colors.navInactive,
        )
    }
}

/** Switch in Cardwise blue, e.g. for settings. */
@Composable
fun CardwiseSwitch(
    checked: Boolean,
    onCheckedChange: (Boolean) -> Unit,
    modifier: Modifier = Modifier,
    width: Dp = 40.dp,
    height: Dp = 24.dp,
) {
    val colors = CardwiseTheme.colors
    val track by animateColorAsState(if (checked) colors.blue else colors.line)
    val thumbSize = height - 6.dp
    val offset = if (checked) width - thumbSize - 3.dp else 3.dp
    Box(
        modifier = modifier
            .size(width, height)
            .clip(CircleShape)
            .background(track)
            .clickable(role = Role.Switch) { onCheckedChange(!checked) },
        contentAlignment = Alignment.CenterStart,
    ) {
        Box(
            Modifier
                .offset(x = offset)
                .size(thumbSize)
                .shadow(1.dp, CircleShape)
                .background(colors.surface, CircleShape),
        )
    }
}
