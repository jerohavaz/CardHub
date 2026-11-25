package cardwise.designsystem.components

import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.aspectRatio
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.offset
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.draw.rotate
import androidx.compose.ui.draw.shadow
import androidx.compose.ui.graphics.Brush
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.style.TextOverflow
import androidx.compose.ui.unit.Dp
import androidx.compose.ui.unit.dp
import cardwise.designsystem.icons.CardwiseIcon
import cardwise.designsystem.icons.IconImage
import cardwise.designsystem.theme.CardwiseTheme

/** Aspect ratio of a real trading card (63 × 88 mm). */
const val CardAspectRatio = 63f / 88f

/** Height of the card image in a [CardTile]. */
val TileArtworkHeight = 195.dp

/**
 * Area for the card image.
 * Shows a placeholder gradient until real card images exist.
 * Later pass [image], e.g. with Coil: `image = { AsyncImage(url, …) }`.
 */
@Composable
fun CardArtwork(
    modifier: Modifier = Modifier,
    placeholderColors: List<Color> = ArtworkPalettes.Lumen,
    cornerRadius: Dp = 8.dp,
    elevation: Dp = 3.dp,
    image: (@Composable () -> Unit)? = null,
) {
    val shape = RoundedCornerShape(cornerRadius)
    Box(
        modifier = modifier
            .aspectRatio(CardAspectRatio)
            .shadow(elevation, shape)
            .clip(shape)
            .background(Brush.linearGradient(placeholderColors)),
        contentAlignment = Alignment.Center,
    ) {
        if (image != null) {
            image()
        } else {
            // Placeholder: hinted card frame
            Box(
                Modifier
                    .fillMaxSize()
                    .padding(6.dp)
                    .border(1.5.dp, Color.White.copy(alpha = 0.55f), RoundedCornerShape(cornerRadius * 0.7f)),
            )
        }
    }
}

/** Placeholder colors per element (until real card images exist). */
object ArtworkPalettes {
    val Lumen = listOf(Color(0xFFDCE9FF), Color(0xFF7FA7F0))
    val Tidal = listOf(Color(0xFFD5F3F1), Color(0xFF5FB8C9))
    val Natur = listOf(Color(0xFFE5F3D2), Color(0xFF8EC27A))
    val Feuer = listOf(Color(0xFFFFE6CC), Color(0xFFE59B6C))
    val Kosmos = listOf(Color(0xFFF2DDF5), Color(0xFFC38BD8))
}

/**
 * Card tile for the Pokédex grid.
 * Tapping opens the details, "+" adds the card to the collection.
 */
@Composable
fun CardTile(
    name: String,
    meta: String,
    price: String,
    rarity: Rarity,
    onClick: () -> Unit,
    onAdd: () -> Unit,
    modifier: Modifier = Modifier,
    artwork: @Composable () -> Unit = { CardArtwork(Modifier.height(TileArtworkHeight)) },
) {
    val colors = CardwiseTheme.colors
    val shape = CardwiseTheme.shapes.card
    Column(
        modifier = modifier
            .shadow(1.dp, shape, ambientColor = ShadowColor, spotColor = ShadowColor)
            .clip(shape)
            .background(colors.surface, shape)
            .border(1.dp, colors.line, shape)
            .clickable(onClickLabel = "Details zu $name", onClick = onClick)
            .padding(start = 9.dp, end = 9.dp, top = 9.dp, bottom = 12.dp),
    ) {
        Box(Modifier.fillMaxWidth(), contentAlignment = Alignment.Center) { artwork() }
        Row(Modifier.fillMaxWidth().padding(top = 8.dp), verticalAlignment = Alignment.CenterVertically) {
            Text(
                name,
                style = CardwiseTheme.typography.labelStrong,
                color = colors.text,
                maxLines = 1,
                overflow = TextOverflow.Ellipsis,
                modifier = Modifier.weight(1f),
            )
            RarityTag(rarity)
        }
        Text(meta, style = CardwiseTheme.typography.badge.copy(fontWeight = null), color = colors.muted, modifier = Modifier.padding(top = 4.dp))
        Row(Modifier.fillMaxWidth().padding(top = 4.dp), verticalAlignment = Alignment.CenterVertically) {
            Text(price, style = CardwiseTheme.typography.label, color = colors.text, modifier = Modifier.weight(1f))
            AddButton(onClick = onAdd, contentDescription = "$name zur Sammlung hinzufügen")
        }
    }
}

/** Collection row: small card, name, price badge, stepper and delete. */
@Composable
fun CollectionItemRow(
    name: String,
    meta: String,
    price: String,
    quantity: Int,
    onQuantityChange: (Int) -> Unit,
    onRemove: () -> Unit,
    onClick: () -> Unit,
    modifier: Modifier = Modifier,
    artwork: @Composable () -> Unit = { CardArtwork(Modifier.width(56.dp), cornerRadius = 4.dp) },
) {
    val colors = CardwiseTheme.colors
    val shape = CardwiseTheme.shapes.xl
    Row(
        modifier = modifier
            .fillMaxWidth()
            .clip(shape)
            .background(colors.surface, shape)
            .border(1.dp, colors.line, shape)
            .clickable(onClick = onClick)
            .padding(11.dp),
        horizontalArrangement = Arrangement.spacedBy(10.dp),
        verticalAlignment = Alignment.CenterVertically,
    ) {
        artwork()
        Column(Modifier.weight(1f)) {
            Text(name, style = CardwiseTheme.typography.labelStrong, color = colors.text)
            Text(meta, style = CardwiseTheme.typography.caption, color = colors.muted, modifier = Modifier.padding(top = 4.dp))
            CardwiseBadge(price, modifier = Modifier.padding(top = 10.dp))
        }
        Column(horizontalAlignment = Alignment.Start, verticalArrangement = Arrangement.spacedBy(8.dp)) {
            Stepper(value = quantity, onValueChange = onQuantityChange)
            Box(
                Modifier
                    .height(32.dp)
                    .width(30.dp)
                    .clip(CardwiseTheme.shapes.sm)
                    .clickable(onClickLabel = "$name entfernen", onClick = onRemove),
                contentAlignment = Alignment.Center,
            ) {
                IconImage(CardwiseIcon.Trash, tint = colors.dangerIcon, size = 16.dp, contentDescription = "$name entfernen")
            }
        }
    }
}

/** Hero banner on the Pokédex start page. */
@Composable
fun HeroBanner(
    eyebrow: String,
    title: String,
    text: String,
    buttonText: String,
    onButtonClick: () -> Unit,
    modifier: Modifier = Modifier,
    artwork: @Composable () -> Unit = {
        CardArtwork(Modifier.width(94.dp).padding(top = 4.dp), cornerRadius = 6.dp, elevation = 8.dp)
    },
) {
    val colors = CardwiseTheme.colors
    val shape = CardwiseTheme.shapes.hero
    Box(
        modifier = modifier
            .fillMaxWidth()
            .clip(shape)
            .background(Brush.linearGradient(listOf(Color(0xFFDCE9FF), Color(0xFFAECBFF))), shape),
    ) {
        // Decorative circles in the top right
        Box(
            Modifier
                .align(Alignment.TopEnd)
                .offset(x = 10.dp, y = (-40).dp)
                .width(145.dp)
                .aspectRatio(1f)
                .border(1.dp, Color.White.copy(alpha = 0.62f), RoundedCornerShape(50)),
        )
        Row(Modifier.padding(17.dp), verticalAlignment = Alignment.CenterVertically) {
            Column(Modifier.weight(1f)) {
                Text(eyebrow.uppercase(), style = CardwiseTheme.typography.badge.copy(letterSpacing = CardwiseTheme.typography.eyebrow.letterSpacing), color = colors.blue)
                Text(title, style = CardwiseTheme.typography.hero, color = colors.heroText, modifier = Modifier.padding(top = 9.dp))
                Text(
                    text,
                    style = CardwiseTheme.typography.caption.copy(lineHeight = CardwiseTheme.typography.caption.fontSize * 1.45),
                    color = colors.heroText.copy(alpha = 0.65f),
                    modifier = Modifier.padding(top = 7.dp, bottom = 12.dp),
                )
                PillButton(buttonText, onClick = onButtonClick)
            }
            Spacer(Modifier.width(12.dp))
            Box(Modifier.rotate(11f)) { artwork() }
        }
    }
}
