package cardwise.designsystem.components

import androidx.compose.animation.animateColorAsState
import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxHeight
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.widthIn
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
import androidx.compose.ui.unit.dp
import cardwise.designsystem.icons.CardwiseIcon
import cardwise.designsystem.icons.IconImage
import cardwise.designsystem.theme.CardwiseTheme

/** The five main tabs of the app. */
enum class MainTab(val label: String, val icon: CardwiseIcon) {
    Pokedex("Pokédex", CardwiseIcon.Home),
    Trading("Trading", CardwiseIcon.Trade),
    Collection("Sammlung", CardwiseIcon.Layers),
    Friends("Freunde", CardwiseIcon.Users),
    Profile("Profil", CardwiseIcon.User),
}

// Figma uses 82 % opacity + blur; without blur we use a higher opacity so it stays readable
private val NavBackground = Color(0xF0F9FBFF)
private val NavBorder = Color(0x6BB4C1DA)
private val NavShadow = Color(0x21364E78)

/**
 * Floating bottom navigation bar.
 * It overlaps the content, so screens need about 110 dp of bottom padding.
 */
@Composable
fun MainNavigationBar(
    selected: MainTab,
    onSelect: (MainTab) -> Unit,
    modifier: Modifier = Modifier,
) {
    val shape = CardwiseTheme.shapes.sheet
    Row(
        modifier = modifier
            .widthIn(max = 420.dp)
            .height(67.dp)
            .shadow(elevation = 14.dp, shape = shape, ambientColor = NavShadow, spotColor = NavShadow)
            .clip(shape)
            .background(NavBackground, shape)
            .border(1.dp, NavBorder, shape)
            .padding(start = 5.dp, end = 5.dp, top = 6.dp, bottom = 7.dp),
        horizontalArrangement = Arrangement.spacedBy(2.dp),
    ) {
        MainTab.entries.forEach { tab ->
            NavItem(tab, selected = tab == selected, onClick = { onSelect(tab) }, modifier = Modifier.weight(1f))
        }
    }
}

@Composable
private fun NavItem(tab: MainTab, selected: Boolean, onClick: () -> Unit, modifier: Modifier) {
    val colors = CardwiseTheme.colors
    val shape = RoundedCornerShape(21.dp)
    val container by animateColorAsState(if (selected) colors.blueSoft else Color.Transparent)
    val content by animateColorAsState(if (selected) colors.blue else colors.navInactive)
    Column(
        modifier = modifier
            .fillMaxHeight()
            .clip(shape)
            .background(container, shape)
            .clickable(role = Role.Tab, onClick = onClick),
        verticalArrangement = Arrangement.spacedBy(3.dp, Alignment.CenterVertically),
        horizontalAlignment = Alignment.CenterHorizontally,
    ) {
        IconImage(tab.icon, tint = content, size = 20.dp)
        Text(tab.label, style = CardwiseTheme.typography.badge.copy(fontWeight = androidx.compose.ui.text.font.FontWeight.Bold), color = content)
    }
}
