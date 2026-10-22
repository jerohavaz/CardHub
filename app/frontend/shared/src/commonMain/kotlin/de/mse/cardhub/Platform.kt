package de.mse.cardhub

interface Platform {
    val name: String
}

expect fun getPlatform(): Platform