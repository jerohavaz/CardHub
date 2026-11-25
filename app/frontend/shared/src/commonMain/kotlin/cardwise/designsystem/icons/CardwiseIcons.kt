package cardwise.designsystem.icons

import androidx.compose.foundation.layout.size
import androidx.compose.material3.Icon
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.unit.Dp
import androidx.compose.ui.unit.dp
import cardwise.shared.generated.resources.Res
import cardwise.shared.generated.resources.ic_camera
import cardwise.shared.generated.resources.ic_chevron_left
import cardwise.shared.generated.resources.ic_chevron_right
import cardwise.shared.generated.resources.ic_close
import cardwise.shared.generated.resources.ic_heart
import cardwise.shared.generated.resources.ic_home
import cardwise.shared.generated.resources.ic_layers
import cardwise.shared.generated.resources.ic_pencil
import cardwise.shared.generated.resources.ic_plus
import cardwise.shared.generated.resources.ic_search
import cardwise.shared.generated.resources.ic_sliders
import cardwise.shared.generated.resources.ic_sparkles
import cardwise.shared.generated.resources.ic_trade
import cardwise.shared.generated.resources.ic_trash
import cardwise.shared.generated.resources.ic_user
import cardwise.shared.generated.resources.ic_users
import org.jetbrains.compose.resources.DrawableResource
import org.jetbrains.compose.resources.painterResource

/**
 * All Icons of the App (Lucide-Stil, from Figma exported and
 * `composeResources/drawable` converted). If new icons need to be added shoudl be done there.
 */
enum class CardwiseIcon(val resource: DrawableResource) {
    Home(Res.drawable.ic_home),
    Trade(Res.drawable.ic_trade),
    Layers(Res.drawable.ic_layers),
    Users(Res.drawable.ic_users),
    User(Res.drawable.ic_user),
    ChevronLeft(Res.drawable.ic_chevron_left),
    ChevronRight(Res.drawable.ic_chevron_right),
    Heart(Res.drawable.ic_heart),
    Plus(Res.drawable.ic_plus),
    Sparkles(Res.drawable.ic_sparkles),
    Search(Res.drawable.ic_search),
    Sliders(Res.drawable.ic_sliders),
    Camera(Res.drawable.ic_camera),
    Trash(Res.drawable.ic_trash),
    Close(Res.drawable.ic_close),
    Pencil(Res.drawable.ic_pencil),
}

@Composable
fun IconImage(
    icon: CardwiseIcon,
    tint: Color,
    modifier: Modifier = Modifier,
    size: Dp = 18.dp,
    contentDescription: String? = null,
) {
    Icon(
        painter = painterResource(icon.resource),
        contentDescription = contentDescription,
        tint = tint,
        modifier = modifier.size(size),
    )
}
