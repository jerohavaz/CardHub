package cardwise.designsystem.components

import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.remember
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Brush
import androidx.compose.ui.graphics.StrokeCap
import androidx.compose.ui.unit.dp
import com.patrykandpatrick.vico.multiplatform.cartesian.CartesianChartHost
import com.patrykandpatrick.vico.multiplatform.cartesian.data.CartesianChartModelProducer
import com.patrykandpatrick.vico.multiplatform.cartesian.data.CartesianLayerRangeProvider
import com.patrykandpatrick.vico.multiplatform.cartesian.data.lineSeries
import com.patrykandpatrick.vico.multiplatform.cartesian.layer.LineCartesianLayer
import com.patrykandpatrick.vico.multiplatform.cartesian.layer.rememberLine
import com.patrykandpatrick.vico.multiplatform.cartesian.layer.rememberLineCartesianLayer
import com.patrykandpatrick.vico.multiplatform.cartesian.rememberCartesianChart
import com.patrykandpatrick.vico.multiplatform.cartesian.rememberVicoScrollState
import com.patrykandpatrick.vico.multiplatform.cartesian.rememberVicoZoomState
import com.patrykandpatrick.vico.multiplatform.common.Fill
import cardwise.designsystem.theme.CardwiseTheme

/**
 * Price history as a line chart with a soft gradient below (Vico).
 * [points] are the prices in chronological order.
 */
@Composable
fun PriceChart(
    points: List<Double>,
    modifier: Modifier = Modifier,
) {
    val colors = CardwiseTheme.colors
    val modelProducer = remember { CartesianChartModelProducer() }
    LaunchedEffect(points) {
        if (points.isNotEmpty()) modelProducer.runTransaction { lineSeries { series(points) } }
    }
    val min = points.minOrNull() ?: 0.0
    val max = points.maxOrNull() ?: 1.0
    val padding = ((max - min) * 0.15).coerceAtLeast(0.5)
    val line = LineCartesianLayer.rememberLine(
        fill = remember(colors.blue) { LineCartesianLayer.LineFill.single(Fill(colors.blue)) },
        stroke = LineCartesianLayer.LineStroke.Continuous(thickness = 2.5.dp, cap = StrokeCap.Round),
        areaFill = remember(colors.blue) {
            LineCartesianLayer.AreaFill.single(
                Fill(Brush.verticalGradient(listOf(colors.blue.copy(alpha = 0.22f), colors.blue.copy(alpha = 0f)))),
            )
        },
        pointConnector = LineCartesianLayer.PointConnector.cubic(),
    )
    CartesianChartHost(
        chart = rememberCartesianChart(
            rememberLineCartesianLayer(
                lineProvider = LineCartesianLayer.LineProvider.series(line),
                rangeProvider = remember(min, max) { CartesianLayerRangeProvider.fixed(minY = min - padding, maxY = max + padding) },
            ),
        ),
        modelProducer = modelProducer,
        modifier = modifier,
        scrollState = rememberVicoScrollState(scrollEnabled = false),
        zoomState = rememberVicoZoomState(zoomEnabled = false),
    )
}

/** Market panel: label, price, trend, chart and month axis. */
@Composable
fun PricePanel(
    label: String,
    price: String,
    trend: String,
    points: List<Double>,
    axisLabels: List<String>,
    modifier: Modifier = Modifier,
    trendPositive: Boolean = true,
) {
    val colors = CardwiseTheme.colors
    Panel(modifier) {
        Row(Modifier.fillMaxWidth(), verticalAlignment = Alignment.Bottom) {
            androidx.compose.foundation.layout.Column(Modifier.weight(1f)) {
                Text(label, style = CardwiseTheme.typography.caption, color = colors.muted)
                Text(price, style = CardwiseTheme.typography.price, color = colors.text)
            }
            Text(
                text = (if (trendPositive) "↗ " else "↘ ") + trend,
                style = CardwiseTheme.typography.label.copy(fontWeight = CardwiseTheme.typography.button.fontWeight),
                color = if (trendPositive) colors.success else colors.danger,
                modifier = Modifier.padding(bottom = 8.dp),
            )
        }
        PriceChart(points, Modifier.fillMaxWidth().height(120.dp).padding(top = 11.dp))
        Row(
            Modifier.fillMaxWidth().padding(horizontal = 3.dp),
            horizontalArrangement = androidx.compose.foundation.layout.Arrangement.SpaceBetween,
        ) {
            axisLabels.forEach { Text(it, style = CardwiseTheme.typography.badge.copy(fontWeight = null), color = colors.muted) }
        }
    }
}
